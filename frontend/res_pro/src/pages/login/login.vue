<template>
    <div class="login-container">
      <div class="login-box">
        <h1>Login</h1>
        <div>
          <label for="username">Username</label>
          <input type="text" id="username" v-model="username" placeholder="Enter username" />
          <label for="password">Password</label>
          <input type="password" id="password" v-model="password" placeholder="Enter password" />
          <button @click="handleLogin">Login</button>
        </div>
        <p class="register-link">
          No account?<span @click="goToRegister">register</span>
        </p>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  const username = ref('');
  const password = ref('');
  const handleLogin = () => {
    console.log('handleLogin triggered');
    if (!username.value || !password.value) {
      uni.showToast({ title: 'Please enter username and password', icon: 'none' });
      return;
    }
  
    uni.request({
      url: 'http://localhost:8080/login',
      method: 'POST',
      data: {
        username: username.value,
        password: password.value
      },
      success: (res) => {
        if (res.statusCode === 200) {
          uni.showToast({ title: 'Login success', icon: 'success' });
          uni.redirectTo({ url: '/pages/Main_page/main' });
        } else {
          uni.showToast({ title: 'Invalid username or password', icon: 'none' });
        }
      },
      fail: () => {
        uni.showToast({ title: 'Request failed', icon: 'none' });
      }
    });
  };
  const goToRegister = () => {
    uni.navigateTo({
      url: '/pages/login/register'
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
  
  .register-link {
    text-align: center;
    font-size: 0.9em;
    margin-top: 10px;
  }
  
  .register-link span {
    color: #007aff;
    cursor: pointer;
    text-decoration: underline;
  }
  </style>