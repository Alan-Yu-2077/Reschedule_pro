<template>
    <div class="login-container">
      <div class="login-box">
        <h1>Register</h1>
        <div>
          <label for="userID">User ID (10 digits)</label>
          <input type="text" id="userID" v-model="userID" placeholder="Enter 10-digit user ID" maxlength="10" />

          <label for="username">Username</label>
          <input type="text" id="username" v-model="username" placeholder="Enter username" />

          <label for="password">Password</label>
          <input type="password" id="password" v-model="password" placeholder="Enter password" />

          <label for="confirmPassword">Confirm Password</label>
          <input type="password" id="confirmPassword" v-model="confirmPassword" placeholder="Confirm password" />

          <button @click="handleRegister">Register</button>
        </div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  const userID = ref('');
  const username = ref('');
  const password = ref('');
  const confirmPassword = ref('');
  
const handleRegister = () => {
  if (!userID.value || !username.value || !password.value || !confirmPassword.value) {
    uni.showToast({ title: 'Please complete all fields', icon: 'none' });
    return;
  }

  // 验证UserID是否为10位数字
  if (!/^\d{10}$/.test(userID.value)) {
    uni.showToast({ title: 'UserID must be 10 digits', icon: 'none' });
    return;
  }

  if (password.value !== confirmPassword.value) {
    uni.showToast({ title: 'Passwords do not match', icon: 'none' });
    return;
  }

  uni.request({
    url: 'http://localhost:8080/register',
    method: 'POST',
    data: {
      userID: userID.value,
      username: username.value,
      password: password.value
    },
    success: (res) => {
      if (res.statusCode === 200) {
        uni.showToast({ title: 'Register success', icon: 'success' });
        uni.navigateBack(); // 返回登录页
      } else {
        uni.showToast({ title: res.data.msg || 'Register failed', icon: 'none' });
      }
    },
    fail: () => {
      uni.showToast({ title: 'Request failed', icon: 'none' });
    }
  });
};
  </script>
  
  <style scoped>
  @import url('https://fonts.googleapis.com/css2?family=Patrick+Hand&display=swap');
  
  .login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
    background: #fdfcfb;
    background-image:
      linear-gradient(to right, #000 1px, transparent 1px),
      linear-gradient(to bottom, #000 1px, transparent 1px);
    background-size: 40px 40px;
  }
  
  .login-box {
    background: rgba(255,255,255,0.8);
    padding: 40px;
    border: 2px dashed #000;
    font-family: 'Patrick Hand', cursive;
    box-shadow: 4px 4px 0px #000;
    max-width: 320px;
    width: 100%;
  }
  
  .login-box h1 {
    text-align: center;
    margin-bottom: 20px;
  }
  
  .login-box label {
    display: block;
    margin-bottom: 5px;
    font-size: 1.1em;
  }
  
  .login-box input {
    width: 100%;
    padding: 8px;
    margin-bottom: 15px;
    border: 2px solid #000;
    outline: none;
    font-family: 'Patrick Hand', cursive;
    font-size: 1em;
  }
  
  .login-box input::placeholder {
    font-family: 'Patrick Hand', cursive;
  }
  
  .login-box button {
    width: 100%;
    padding: 10px;
    border: 2px solid #000;
    background: #fff;
    cursor: pointer;
    font-family: 'Patrick Hand', cursive;
    font-size: 1.2em;
    transition: transform 0.1s;
  }
  
  .login-box button:active {
    transform: translate(2px, 2px);
  }
  </style>