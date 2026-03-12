import pytest
import requests
import jwt
import base64
import json

BASE_URL = "http://localhost:9091/api/v1"
EMAIL = "admin@demo.com"
PASSWORD = "password123"

@pytest.fixture
def login():
    """Helper fixture to get a fresh login token before each test needing it"""
    res = requests.post(f"{BASE_URL}/auth/login", json={"email": EMAIL, "password": PASSWORD})
    assert res.status_code == 200, f"Login failed: {res.text}"
    return res.json()

# TC-01
def test_login_success():
    """Login thành công, trả về 200, có access_token, refresh_token"""
    res = requests.post(f"{BASE_URL}/auth/login", json={"email": EMAIL, "password": PASSWORD})
    assert res.status_code == 200
    data = res.json()
    assert "access_token" in data
    assert "refresh_token" in data
    assert "token_type" in data

# TC-02
def test_jwt_structure(login):
    """Token phải có 3 phần, Payload phải decode được giải mã Base64 và có jti"""
    token = login["access_token"]
    parts = token.split(".")
    assert len(parts) == 3, "JWT must have 3 parts"
    
    # Add padding to base64
    payload_b64 = parts[1]
    payload_b64 += "=" * ((4 - len(payload_b64) % 4) % 4)
    payload_bytes = base64.b64decode(payload_b64)
    payload = json.loads(payload_bytes)
    
    assert "jti" in payload
    assert "session_id" in payload
    assert payload["email"] == EMAIL

# TC-03
def test_jti_uniqueness():
    """Tính duy nhất của JTI: đăng nhập 2 lần, 2 token phải khác nhau"""
    res1 = requests.post(f"{BASE_URL}/auth/login", json={"email": EMAIL, "password": PASSWORD}).json()
    res2 = requests.post(f"{BASE_URL}/auth/login", json={"email": EMAIL, "password": PASSWORD}).json()
    
    def get_jti(token):
        parts = token.split(".")
        payload_b64 = parts[1] + "=" * ((4 - len(parts[1]) % 4) % 4)
        return json.loads(base64.b64decode(payload_b64))["jti"]
        
    assert get_jti(res1["access_token"]) != get_jti(res2["access_token"])

# TC-04
def test_replay_attack(login):
    """Chống Replay Attack: dùng token sau khi logout (bị revoke) -> 401"""
    token = login["access_token"]
    headers = {"Authorization": f"Bearer {token}"}
    
    # 1. Access OK
    res1 = requests.get(f"{BASE_URL}/protected/health", headers=headers)
    assert res1.status_code == 200
    
    # 2. Server revokes token via logout
    res_logout = requests.post(f"{BASE_URL}/auth/logout", headers=headers)
    assert res_logout.status_code == 200
    
    # 3. Replay attack: try accessing the route again
    res3 = requests.get(f"{BASE_URL}/protected/health", headers=headers)
    assert res3.status_code == 401, "Replayed token should be rejected"

# TC-05
def test_logout_and_revoke(login):
    """TC-05: Logout và Auth thu hồi JTI thành công, Redis reject"""
    import time
    time.sleep(1) # sleep to avoid rate limiting
    token = login["access_token"]
    headers = {"Authorization": f"Bearer {token}"}
    
    # 1. Access protected route
    res1 = requests.get(f"{BASE_URL}/protected/health", headers=headers)
    assert res1.status_code == 200
    
    # 2. Logout
    res_logout = requests.post(f"{BASE_URL}/auth/logout", headers=headers)
    assert res_logout.status_code == 200
    
    # 3. Double logout is idempotent (returns 200) in current implementation
    res_double = requests.post(f"{BASE_URL}/auth/logout", headers=headers)
    assert res_double.status_code in [200, 401]
    
    # 4. Access protected resource again -> Must be 401
    res_access_after_logout = requests.get(f"{BASE_URL}/protected/health", headers=headers)
    assert res_access_after_logout.status_code == 401

# TC-06
def test_tampered_payload(login):
    """Thay đổi nội dung payload (sửa role -> superadmin)"""
    import time
    time.sleep(1)
    token = login["access_token"]
    parts = token.split(".")
    
    payload_b64 = parts[1] + "=" * ((4 - len(parts[1]) % 4) % 4)
    payload = json.loads(base64.b64decode(payload_b64))
    
    # Hacker modifies role
    payload["role"] = "hacker"
    
    # Encode back to Base64 (remove '=' padding)
    tampered_payload_b64 = base64.b64encode(json.dumps(payload).encode()).decode().rstrip("=")
    tampered_token = f"{parts[0]}.{tampered_payload_b64}.{parts[2]}"
    
    headers = {"Authorization": f"Bearer {tampered_token}"}
    res = requests.get(f"{BASE_URL}/protected/health", headers=headers)
    assert res.status_code == 401, "Tampered signature should fail verification"

# TC-07
def test_fake_jti():
    """Tạo JWT giả mạo bằng jwt.encode với self-signed secret"""
    import time
    time.sleep(1)
    # RSA is required anyway, but let's try bypass with HS256 as some weak implementations allow
    fake_token = jwt.encode(
        {"user_id": "80882efa-b305-42ad-ad5a-a19cda59590a", "jti": "fake-jti-12345", "role": "admin"}, 
        "fake-secret-key", 
        algorithm="HS256"
    )
    headers = {"Authorization": f"Bearer {fake_token}"}
    res = requests.get(f"{BASE_URL}/protected/health", headers=headers)
    assert res.status_code == 401, "Fake token should fail verification constraint"
