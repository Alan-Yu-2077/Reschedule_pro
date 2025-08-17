<template>
  <div class="admin-container">
    <div class="admin-header">
      <h1>Admin Dashboard</h1>
      <div class="header-actions">
        <button class="refresh-btn" @click="refreshAllData">🔄 Refresh</button>
        <button class="logout-btn" @click="handleLogout">Logout</button>
      </div>
    </div>

    <!-- Tab导航 -->
    <div class="tab-navigation">
      <button 
        v-for="tab in tabs" 
        :key="tab.key"
        class="tab-btn"
        :class="{ active: activeTab === tab.key }"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- 数据展示区域 -->
    <div class="data-content">
      <!-- 用户数据 -->
      <div v-if="activeTab === 'users'" class="data-section">
        <h2>Users ({{ users.length }})</h2>
        <div class="user-actions">
          <button class="add-user-btn" @click="showAddUserModal = true">Add User</button>
        </div>
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>User ID</th>
                <th>Username</th>
                <th>Password</th>
                <th>User Type</th>
                <th>Created At</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in users" :key="user.userID">
                <td>
                  <span v-if="!user.isEditing">{{ user.userID }}</span>
                  <input 
                    v-else 
                    v-model="user.newUserID" 
                    type="text" 
                    class="password-input"
                    placeholder="10-digit ID"
                    maxlength="10"
                  />
                </td>
                <td>
                  <span v-if="!user.isEditing">{{ user.username }}</span>
                  <input 
                    v-else 
                    v-model="user.newUsername" 
                    type="text" 
                    class="password-input"
                    placeholder="Username"
                  />
                </td>
                <td>
                  <span v-if="!user.isEditing">{{ user.password }}</span>
                  <input 
                    v-else 
                    v-model="user.newPassword" 
                    type="text" 
                    class="password-input"
                    placeholder="Enter new password"
                  />
                </td>
                <td>
                  <span v-if="!user.isEditing">{{ user.userType || 'user' }}</span>
                  <select v-else v-model="user.newUserType" class="password-input">
                    <option value="viewer">viewer</option>
                    <option value="user">user</option>
                    <option value="teacher">teacher</option>
                    <option value="admin">admin</option>
                  </select>
                </td>
                <td>{{ formatDate(user.createdAt) }}</td>
                <td>
                  <button 
                    v-if="!user.isEditing" 
                    class="edit-btn" 
                    @click="startEditPassword(user)"
                  >
                    Edit
                  </button>
                  <div v-else class="edit-actions">
                    <button class="save-btn" @click="savePassword(user)">Save</button>
                    <button class="cancel-btn" @click="cancelEditPassword(user)">Cancel</button>
                  </div>
                  <button class="delete-btn" @click="confirmDeleteUser(user)">Delete</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 班级数据 -->
      <div v-if="activeTab === 'classes'" class="data-section">
        <h2>Classes ({{ classes.length }})</h2>
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Class Name</th>
                <th>Created At</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="cls in classes" :key="cls.ID">
                <td>{{ cls.ID }}</td>
                <td>{{ cls.name }}</td>
                <td>{{ formatDate(cls.CreatedAt) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 课程表数据 -->
      <div v-if="activeTab === 'schedules'" class="data-section">
        <h2>Schedules ({{ schedules.length }})</h2>
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Class</th>
                <th>Course</th>
                <th>Week</th>
                <th>Time Slot</th>
                <th>Created At</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="schedule in schedules" :key="schedule.ID">
                <td>{{ schedule.ID }}</td>
                <td>{{ schedule.class?.name || 'N/A' }}</td>
                <td>{{ schedule.course?.name || 'N/A' }}</td>
                <td>{{ schedule.weekNumber }}</td>
                <td>{{ getTimeSlotDisplay(schedule.timeSlotRow, schedule.timeSlotCol) }}</td>
                <td>{{ formatDate(schedule.CreatedAt) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 日志数据 -->
      <div v-if="activeTab === 'logs'" class="data-section">
        <h2>Activity Logs ({{ logs.length }})</h2>
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Message</th>
                <th>Created At</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="log in logs" :key="log.ID">
                <td>{{ log.ID }}</td>
                <td>{{ log.message }}</td>
                <td>{{ formatDate(log.CreatedAt) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Add User Modal -->
    <div v-if="showAddUserModal" class="modal-overlay" @click.self="showAddUserModal = false">
      <div class="modal-content add-user-modal" @click.stop>
        <h3>Add New User</h3>
        <div class="input-group">
          <label>User ID (10 digits):</label>
          <input v-model="newUser.userID" type="text" class="user-input" placeholder="10-digit ID" maxlength="10" />
        </div>
        <div class="input-group">
          <label>Username:</label>
          <input v-model="newUser.username" type="text" class="user-input" placeholder="Username" />
        </div>
        <div class="input-group">
          <label>Password:</label>
          <input v-model="newUser.password" type="text" class="user-input" placeholder="Password" />
        </div>
        <div class="input-group">
          <label>User Type:</label>
          <select v-model="newUser.userType" class="user-input">
            <option value="viewer">viewer</option>
            <option value="user">user</option>
            <option value="teacher">teacher</option>
            <option value="admin">admin</option>
          </select>
        </div>
        <div class="modal-actions">
          <button class="save-btn" @click="addUser">Add</button>
          <button class="cancel-btn" @click="showAddUserModal = false">Cancel</button>
        </div>
      </div>
    </div>

    <!-- Delete User Confirmation Modal -->
    <div v-if="showDeleteUserModal" class="modal-overlay" @click.self="showDeleteUserModal = false">
      <div class="modal-content delete-user-modal" @click.stop>
        <h3>Confirm Delete</h3>
        <p>Are you sure you want to delete user <b>{{ userToDelete?.username }}</b>?</p>
        <div class="modal-actions">
          <button class="delete-btn" @click="deleteUser">Confirm Delete</button>
          <button class="cancel-btn" @click="showDeleteUserModal = false">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue';

// 响应式数据
const activeTab = ref('users');
const users = ref([]);
const classes = ref([]);
const schedules = ref([]);
const logs = ref([]);

// Tab配置
const tabs = [
  { key: 'users', label: 'Users' },
  { key: 'classes', label: 'Classes' },
  { key: 'schedules', label: 'Schedules' },
  { key: 'logs', label: 'Logs' }
];

// 添加/删除用户相关
const showAddUserModal = ref(false);
const showDeleteUserModal = ref(false);
const newUser = reactive({ userID: '', username: '', password: '', userType: 'viewer' });
const userToDelete = ref(null);

// 加载用户数据
const loadUsers = async () => {
  try {
    const response = await uni.request({
      url: 'http://localhost:8080/admin/users',
      method: 'GET'
    });
    if (response.statusCode === 200) {
      users.value = response.data.users || [];
    }
  } catch (error) {
    console.error('Failed to load users:', error);
    uni.showToast({ title: 'Failed to load users', icon: 'none' });
  }
};

// 加载班级数据
const loadClasses = async () => {
  try {
    const response = await uni.request({
      url: 'http://localhost:8080/admin/classes',
      method: 'GET'
    });
    if (response.statusCode === 200) {
      classes.value = response.data.classes || [];
    }
  } catch (error) {
    console.error('Failed to load classes:', error);
    uni.showToast({ title: 'Failed to load classes', icon: 'none' });
  }
};

// 加载课程表数据
const loadSchedules = async () => {
  try {
    const response = await uni.request({
      url: 'http://localhost:8080/admin/schedules',
      method: 'GET'
    });
    if (response.statusCode === 200) {
      schedules.value = response.data.schedules || [];
    }
  } catch (error) {
    console.error('Failed to load schedules:', error);
    uni.showToast({ title: 'Failed to load schedules', icon: 'none' });
  }
};

// 加载日志数据
const loadLogs = async () => {
  try {
    const response = await uni.request({
      url: 'http://localhost:8080/admin/logs',
      method: 'GET'
    });
    if (response.statusCode === 200) {
      logs.value = response.data.logs || [];
    }
  } catch (error) {
    console.error('Failed to load logs:', error);
    uni.showToast({ title: 'Failed to load logs', icon: 'none' });
  }
};

// 刷新所有数据
const refreshAllData = async () => {
  uni.showToast({ title: 'Refreshing data...', icon: 'loading' });
  await Promise.all([
    loadUsers(),
    loadClasses(),
    loadSchedules(),
    loadLogs()
  ]);
  uni.showToast({ title: 'Data refreshed', icon: 'success' });
};

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return 'N/A';
  const date = new Date(dateString);
  return date.toLocaleString();
};

// 获取时间段显示
const getTimeSlotDisplay = (row, col) => {
  const periods = ['Morning Slot 1', 'Morning Slot 2', 'Afternoon Slot 1', 'Afternoon Slot 2', 'Evening Slot'];
  const days = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
  return `${periods[row] || 'Unknown'} - ${days[col] || 'Unknown'}`;
};

// 开始编辑用户信息
const startEditPassword = (user) => {
  user.isEditing = true;
  user.newUserID = user.userID; // 初始值为当前用户ID
  user.newUsername = user.username; // 初始值为当前用户名
  user.newPassword = user.password; // 初始值为当前密码
  user.newUserType = user.userType || 'viewer'; // 初始值为当前用户类型
};

// 保存用户信息
const savePassword = async (user) => {
  if (!user.newUserID || !user.newUsername || !user.newPassword || !user.newUserType) {
    uni.showToast({ title: 'All fields cannot be empty', icon: 'none' });
    return;
  }

  // 验证UserID是否为10位数字
  if (!/^\d{10}$/.test(user.newUserID)) {
    uni.showToast({ title: 'User ID must be 10 digits', icon: 'none' });
    return;
  }

  try {
    const response = await uni.request({
      url: `http://localhost:8080/admin/users/${user.userID}`,
      method: 'PUT',
      data: {
        newUserID: user.newUserID,
        newUsername: user.newUsername,
        newPassword: user.newPassword,
        newUserType: user.newUserType
      }
    });

    if (response.statusCode === 200) {
      // 更新本地数据
      user.userID = user.newUserID;
      user.username = user.newUsername;
      user.password = user.newPassword;
      user.userType = user.newUserType;
      user.isEditing = false;
      
      // 清理临时数据
      delete user.newUserID;
      delete user.newUsername;
      delete user.newPassword;
      delete user.newUserType;
      
      uni.showToast({ title: 'User information updated successfully', icon: 'success' });
    } else {
      uni.showToast({ title: response.data.error || 'Update failed', icon: 'none' });
    }
  } catch (error) {
    console.error('Failed to update user:', error);
    uni.showToast({ title: 'Update failed', icon: 'none' });
  }
};

// 取消编辑用户信息
const cancelEditPassword = (user) => {
  user.isEditing = false;
  delete user.newUserID;
  delete user.newUsername;
  delete user.newPassword;
  delete user.newUserType;
};

// 登出
const handleLogout = () => {
  uni.clearStorageSync(); // 清除本地存储
  uni.showToast({ title: 'Logged out successfully', icon: 'success' });
  setTimeout(() => {
    uni.reLaunch({
      url: '/pages/login/login'
    });
  }, 1000);
};

// 添加用户
const addUser = async () => {
  if (!newUser.userID.trim() || !newUser.username.trim() || !newUser.password.trim()) {
    uni.showToast({ title: 'User ID, username and password cannot be empty', icon: 'none' });
    return;
  }
  
  // 验证UserID是否为10位数字
  if (!/^\d{10}$/.test(newUser.userID)) {
    uni.showToast({ title: 'User ID must be 10 digits', icon: 'none' });
    return;
  }
  
  try {
    const response = await uni.request({
      url: 'http://localhost:8080/admin/users',
      method: 'POST',
      data: {
        userID: newUser.userID,
        username: newUser.username,
        password: newUser.password,
        userType: newUser.userType
      }
    });
    if (response.statusCode === 200) {
      showAddUserModal.value = false;
      newUser.userID = '';
      newUser.username = '';
      newUser.password = '';
      newUser.userType = 'viewer';
      await loadUsers();
      uni.showToast({ title: 'User added successfully', icon: 'success' });
    } else {
      uni.showToast({ title: response.data.error || 'Add failed', icon: 'none' });
    }
  } catch (error) {
    uni.showToast({ title: 'Add failed', icon: 'none' });
  }
};

// 确认删除用户
const confirmDeleteUser = (user) => {
  userToDelete.value = user;
  showDeleteUserModal.value = true;
};

// 删除用户
const deleteUser = async () => {
  if (!userToDelete.value) return;
  try {
    const response = await uni.request({
      url: `http://localhost:8080/admin/users/${userToDelete.value.userID}`,
      method: 'DELETE'
    });
    if (response.statusCode === 200) {
      showDeleteUserModal.value = false;
      userToDelete.value = null;
      await loadUsers();
      uni.showToast({ title: 'User deleted successfully', icon: 'success' });
    } else {
      uni.showToast({ title: response.data.error || 'Delete failed', icon: 'none' });
    }
  } catch (error) {
    uni.showToast({ title: 'Delete failed', icon: 'none' });
  }
};

// 页面加载时获取数据
onMounted(() => {
  refreshAllData();
});
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Patrick+Hand&display=swap');

.admin-container {
  min-height: 100vh;
  background: #fdfcfb;
  font-family: 'Patrick Hand', cursive;
  padding: 20px;
}

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 2px dashed #000;
  margin-bottom: 30px;
}

.admin-header h1 {
  margin: 0;
  font-size: 2em;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.refresh-btn, .logout-btn {
  padding: 10px 20px;
  border: 2px solid #000;
  background: #fff;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 1.1em;
  transition: transform 0.1s;
}

.refresh-btn:hover, .logout-btn:hover {
  transform: translate(2px, 2px);
}

.tab-navigation {
  display: flex;
  gap: 10px;
  margin-bottom: 30px;
  border-bottom: 2px dashed #000;
  padding-bottom: 10px;
}

.tab-btn {
  padding: 12px 24px;
  border: 2px solid #000;
  background: #fff;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 1.1em;
  transition: all 0.2s;
}

.tab-btn.active {
  background: #000;
  color: #fff;
}

.tab-btn:hover {
  transform: translateY(-2px);
}

.data-content {
  min-height: 400px;
}

.data-section h2 {
  margin-bottom: 20px;
  color: #333;
  font-size: 1.5em;
}

.table-container {
  overflow-x: auto;
  border: 2px solid #000;
  border-radius: 8px;
  background: white;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9em;
}

.data-table th,
.data-table td {
  padding: 12px 8px;
  text-align: left;
  border-bottom: 1px solid #ddd;
}

.data-table th {
  background: #f5f5f5;
  font-weight: bold;
  border-bottom: 2px solid #000;
}

.data-table tr:hover {
  background: #f9f9f9;
}

.data-table td {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.password-input {
  width: 100%;
  padding: 4px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: 'Patrick Hand', cursive;
  font-size: 0.9em;
}

.edit-btn, .save-btn, .cancel-btn {
  padding: 4px 8px;
  border: 1px solid #000;
  background: white;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 0.8em;
  margin: 0 2px;
  transition: all 0.2s;
}

.edit-btn:hover, .save-btn:hover, .cancel-btn:hover {
  transform: translateY(-1px);
}

.save-btn {
  background: #4CAF50;
  color: white;
  border-color: #4CAF50;
}

.cancel-btn {
  background: #f44336;
  color: white;
  border-color: #f44336;
}

.edit-actions {
  display: flex;
  gap: 4px;
}

.user-actions {
  margin-bottom: 10px;
}

.add-user-btn {
  padding: 8px 18px;
  border: 2px solid #000;
  background: #fff;
  cursor: pointer;
  font-family: 'Patrick Hand', cursive;
  font-size: 1.1em;
  border-radius: 6px;
  margin-bottom: 10px;
  transition: all 0.2s;
}

.add-user-btn:hover {
  background: #e3f2fd;
}

.delete-btn {
  background: #f44336;
  color: white;
  border-color: #f44336;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.8em;
  margin-left: 2px;
}

.delete-btn:hover {
  background: #d32f2f;
}

.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  padding: 30px 30px 20px 30px;
  min-width: 320px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.2);
}

.add-user-modal h3, .delete-user-modal h3 {
  margin-top: 0;
  margin-bottom: 20px;
}

.input-group {
  margin-bottom: 18px;
}

.user-input {
  width: 100%;
  padding: 8px;
  border: 1.5px solid #ddd;
  border-radius: 6px;
  font-size: 1em;
  font-family: 'Patrick Hand', cursive;
}

.modal-actions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
  justify-content: flex-end;
}
</style> 