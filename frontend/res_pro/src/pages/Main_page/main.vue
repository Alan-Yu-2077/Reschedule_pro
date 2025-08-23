<template>
  <div class="main-container">
    <!-- 侧边栏切换按钮 -->
    <div class="sidebar-toggle" @click="toggleSidebar" title="Toggle Sidebar">
      <div class="hamburger-icon">
        <span></span>
        <span></span>
        <span></span>
      </div>
      <span class="toggle-text">Menu</span>
    </div>

    <!-- 背景遮罩 -->
    <div class="sidebar-overlay" :class="{ active: sidebarExpanded }" @click="closeSidebar"></div>

    <!-- 侧边栏 -->
    <div class="sidebar" :class="{ 'sidebar-expanded': sidebarExpanded }">
      <div class="class-header">
        <h2>Class List</h2>
        <!-- 根据用户类型显示不同的加号功能 -->
        <button v-if="userType === 'user'" class="add-btn" @click="goImportPage">➕</button>
        <button v-else-if="userType === 'teacher'" class="add-btn" @click="openTeacherPlanModal">➕</button>
      </div>
      <ul>
        <li v-for="cls in classList" :key="cls" @click="selectClass(cls)">
          {{ cls }}
        </li>
      </ul>

    </div>

    <div class="schedule-panel">
      <div class="schedule-top">
        <div class="controls">
          <button @click="prevWeek">←</button>
        </div>
        <div class="title">
          <h2>{{ currentClass }} - Week {{ currentWeek }} Schedule</h2>
        </div>
        <div class="controls">
          <button @click="nextWeek">→</button>
          <button v-if="userType === 'teacher'" @click="openTeacherPlanViewer">View Plan</button>
        </div>
      </div>
      <div class="schedule-body">
        <table class="schedule-table">
          <thead>
            <tr>
              <th>Period</th>
              <th v-for="day in ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']" :key="day">{{ day }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(period, i) in ['Morning Slot 1', 'Morning Slot 2', 'Afternoon Slot 1', 'Afternoon Slot 2', 'Evening Slot']" :key="i">
              <td>{{ userType === 'teacher' ? getPeriodTime(i) : period }}</td>
              <td v-for="(day, j) in 7" :key="j">
                <button 
                  class="slot-btn" 
                  :class="{ 
                    'selected': selectedCourse && selectedCourse.row === i && selectedCourse.col === j,
                    'target': isSelectingTarget && selectedCourse && selectedCourse.row === i && selectedCourse.col === j,
                    'busy': isSlotBusy(i, j)
                  }"
                  @click="handleSlotClick(i, j)"
                >
                  {{ getCourseName(i, j) }}
                </button>
                <div 
                  v-if="tooltipSlot && tooltipSlot.row === i && tooltipSlot.col === j" 
                  class="slot-tooltip"
                >
                  {{ getSlotTime(i, j) }}
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="course-swap-panel" v-if="selectedCourse">
        <h3>Course Swap</h3>
        <div class="swap-info">
          <p>Selected: {{ getCourseName(selectedCourse.row, selectedCourse.col) }} (Week {{ selectedCourse.sourceWeek }}, {{ getDayName(selectedCourse.col) }}, {{ getPeriodName(selectedCourse.row) }})</p>
          <p>Current Week: {{ currentWeek }} - Click on a slot below to move the course here</p>
          <div class="swap-actions">
            <button @click="cancelSelection" class="cancel-btn">Cancel</button>
            <button @click="deleteCourse" class="delete-btn">Delete Course</button>
            <button @click="startTargetSelection" class="swap-btn">Move to Current Week</button>
          </div>
        </div>
      </div>
      
      <div class="target-selection-panel" v-if="isSelectingTarget">
        <h3>Select Target Position</h3>
        <div class="target-info">
          <p>Click on a slot to move "{{ getCourseName(selectedCourse.row, selectedCourse.col) }}" there</p>
          <button @click="cancelTargetSelection" class="cancel-btn">Cancel</button>
        </div>
      </div>

      <!-- 添加新课程面板 -->
      <div class="add-course-panel" v-if="isAddingNewCourse">
        <h3>Add New Course</h3>
        <div class="add-course-info">
          <p>Adding course to: Week {{ newCourseSlot.week }}, {{ getDayName(newCourseSlot.col) }}, {{ getPeriodName(newCourseSlot.row) }}</p>
          <div class="add-course-form">
            <input 
              v-model="newCourseName" 
              type="text" 
              placeholder="Enter course name" 
              class="course-name-input"
              maxlength="30"
              @keyup.enter="addNewCourse(newCourseName)"
            />
            <div class="add-course-actions">
              <button @click="cancelAddNewCourse" class="cancel-btn">Cancel</button>
              <button @click="addNewCourse(newCourseName)" class="confirm-add-btn">Add Course</button>
            </div>
          </div>
        </div>
      </div>

      <!-- 教师计划录入弹窗 -->
      <div class="teacher-modal" v-if="showTeacherPlanModal">
        <div class="teacher-modal-content">
          <div class="teacher-plan-header">
            <h4>Teacher Plan</h4>
          </div>
          <div class="teacher-plan-body">
            <div class="input-group">
              <label>Class ID</label>
              <input class="plan-input" v-model="planClassName" placeholder="Enter class ID/name" />
            </div>

            <div class="week-type-selection">
              <label>Week Type</label>
              <div class="week-type-buttons">
                <button 
                  class="week-type-btn" 
                  :class="{ active: planWeekType === 'continuous' }"
                  @click="planWeekType = 'continuous'"
                >Continuous Weeks</button>
                <button 
                  class="week-type-btn" 
                  :class="{ active: planWeekType === 'discrete' }"
                  @click="planWeekType = 'discrete'"
                >Discrete Weeks</button>
              </div>
            </div>

            <div v-if="planWeekType === 'continuous'" class="continuous-weeks">
              <div class="week-range">
                <div class="range-input">
                  <label>Start Week</label>
                  <select v-model.number="planStartWeek" class="week-select">
                    <option v-for="w in 20" :key="'sw'+w" :value="w">{{ w }}</option>
                  </select>
                </div>
                <div class="range-input">
                  <label>End Week</label>
                  <select v-model.number="planEndWeek" class="week-select">
                    <option v-for="w in 20" :key="'ew'+w" :value="w">{{ w }}</option>
                  </select>
                </div>
              </div>
            </div>

            <div v-else class="discrete-weeks">
              <label>Select Weeks</label>
              <div class="week-grid">
                <button 
                  v-for="w in 20" :key="'dw'+w" 
                  class="week-number-btn"
                  :class="{ selected: planSelectedWeeks.includes(w) }"
                  @click="togglePlanWeek(w)"
                >{{ w }}</button>
              </div>
            </div>

            <div class="plan-grid-wrapper">
              <div class="plan-grid-title">Select Teaching Slots</div>
              <table class="plan-table">
                <thead>
                  <tr>
                    <th>Period</th>
                    <th v-for="day in ['Mon','Tue','Wed','Thu','Fri','Sat','Sun']" :key="day">{{ day }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(period, i) in ['Morning Slot 1','Morning Slot 2','Afternoon Slot 1','Afternoon Slot 2','Evening Slot']" :key="i">
                    <td>{{ getPeriodTime(i) }}</td>
                    <td v-for="j in 7" :key="j">
                      <button class="slot-btn plan-slot"
                        :class="{ 'plan-selected': isPlanSelected(i, j-1) }"
                        @click="togglePlanSlot(i, j-1)">
                        {{ isPlanSelected(i, j-1) ? '✓' : '' }}
                      </button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <div class="teacher-plan-footer">
            <button class="close-btn" @click="closeTeacherPlanModal">Cancel</button>
            <button class="save-btn" @click="saveTeacherPlan">Save</button>
            <button class="clear-btn" @click="clearAllTeacherPlans">Clear All My Plans</button>
          </div>
        </div>
      </div>

      <!-- 教师计划查看弹窗（只读） -->
      <div class="teacher-modal" v-if="showTeacherPlanViewer">
        <div class="teacher-modal-content">
          <div class="teacher-plan-header">
            <h4>My Weekly Plan</h4>
          </div>
          <div class="teacher-plan-body">
            <div class="week-type-selection">
              <label>Week</label>
              <div class="viewer-week-switch">
                <button class="week-nav" @click="prevViewerWeek">←</button>
                <span class="week-label">Week {{ viewerWeek }}</span>
                <button class="week-nav" @click="nextViewerWeek">→</button>
              </div>
            </div>

            <div class="plan-grid-wrapper">
              <div class="plan-grid-title">Busy Slots</div>
              <table class="plan-table">
                <thead>
                  <tr>
                    <th>Period</th>
                    <th v-for="day in ['Mon','Tue','Wed','Thu','Fri','Sat','Sun']" :key="day">{{ day }}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(period, i) in ['Morning Slot 1','Morning Slot 2','Afternoon Slot 1','Afternoon Slot 2','Evening Slot']" :key="i">
                    <td>{{ getPeriodTime(i) }}</td>
                    <td v-for="j in 7" :key="j">
                      <div class="viewer-slot" :class="{ 'busy': isViewerBusy(i, j-1) }">
                        <span class="dot" v-if="isViewerBusy(i, j-1)"></span>
                        <span class="cls-text" v-if="isViewerBusy(i, j-1)">{{ getViewerClassText(i, j-1) }}</span>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <div class="teacher-plan-footer">
            <button class="close-btn" @click="closeTeacherPlanViewer">Close</button>
          </div>
        </div>
      </div>

      <div class="log-panel">
        <h3>Activity Log</h3>
        <div class="log-content">
          <div v-for="(log, index) in logs" :key="index">{{ log }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { onLoad } from '@dcloudio/uni-app';

const classList = ref([]);
const currentClass = ref('');
const logs = ref([]);
const scheduleData = ref([]);
const selectedCourse = ref(null); // 选中的课程
const isSelectingTarget = ref(false); // 是否正在选择目标位置
const isAddingNewCourse = ref(false); // 是否正在添加新课程
const newCourseSlot = ref(null); // 新课程的位置信息
const newCourseName = ref(''); // 新课程名称
const tooltipSlot = ref(null); // 当前显示时间提示的格子 {row,col}
const sidebarExpanded = ref(false); // 侧边栏默认收起
const userType = ref(''); // 用户类型
const showTeacherPlanModal = ref(false); // 教师计划弹窗
const showTeacherPlanViewer = ref(false); // 教师计划查看弹窗
const planClassName = ref('');
const planWeekType = ref('continuous');
const planStartWeek = ref(1);
const planEndWeek = ref(1);
const planSelectedWeeks = ref([]);
const planSelectedSlots = ref(new Set()); // key: `${row}-${col}`
const busySlots = ref(new Set()); // 当前周的忙碌集合，key: `${row}-${col}`
const userID = ref('');
const nameToId = ref({});

// 加载班级列表
const loadClasses = async () => {
  try {
    const response = await uni.request({
      url: 'http://localhost:8080/api/schedule/classes',
      method: 'GET',
      header: {
        'Content-Type': 'application/json'
      }
    });

    console.log('Classes API response:', response);

    if (response.statusCode === 200) {
      const classes = response.data.classes || [];
      console.log('Raw classes data:', classes);
      
      classList.value = classes.map(cls => cls.name);
      const map = {};
      classes.forEach(c => { 
        const id = c.ID ?? c.id ?? c.Id;
        map[c.name] = id;
        console.log(`Mapping class "${c.name}" to ID:`, id);
      });
      nameToId.value = map;
      console.log('Final nameToId mapping:', nameToId.value);
      console.log('Classes loaded:', classList.value);
      
      // 如果有班级，选择第一个
      if (classList.value.length > 0 && !currentClass.value) {
        currentClass.value = classList.value[0];
        console.log('Set current class to:', currentClass.value);
        loadSchedule();
        loadLogs();
      }
    } else {
      console.error('Failed to load classes, status:', response.statusCode);
      uni.showToast({ title: 'Failed to load classes', icon: 'none' });
    }
  } catch (error) {
    console.error('Failed to load classes:', error);
    uni.showToast({ title: 'Failed to load classes', icon: 'none' });
  }
};

const selectClass = (cls) => {
  console.log('Selecting class:', cls);
  currentClass.value = cls;
  logs.value.push(`Switched to class: ${cls}`);
  loadSchedule();
  loadLogs();
  // 选择班级后自动收起侧边栏
  sidebarExpanded.value = false;
  loadBusySlots();
};

const router = useRouter();
const goImportPage = () => {
  uni.navigateTo({
    url: '/pages/Main_page/register_class'
  });
};

const currentWeek = ref(1);
const prevWeek = () => {
  if (currentWeek.value > 1) currentWeek.value--;
  loadSchedule(); // 重新加载当前周的课表
  loadBusySlots();
};
const nextWeek = () => {
  if (currentWeek.value < 20) currentWeek.value++;
  loadSchedule(); // 重新加载当前周的课表
  loadBusySlots();
};

// 加载课表数据
const loadSchedule = async () => {
  if (!currentClass.value) return;
  
  try {
    const response = await uni.request({
      url: `http://localhost:8080/api/schedule/class/${encodeURIComponent(currentClass.value)}/week/${currentWeek.value}`,
      method: 'GET'
    });

    console.log('Schedule API response:', response);

    if (response.statusCode === 200) {
      scheduleData.value = response.data.schedules || [];
      console.log('Schedule loaded for week', currentWeek.value, ':', scheduleData.value);
    } else {
      console.error('Failed to load schedule, status:', response.statusCode);
      uni.showToast({ title: 'Failed to load schedule', icon: 'none' });
    }
  } catch (error) {
    console.error('Failed to load schedule:', error);
      uni.showToast({ title: 'Failed to load schedule', icon: 'none' });
    }
};

// 班级日志（从 admin/logs 获取后前端过滤）
const loadLogs = async () => {
  if (!currentClass.value) return;
  try {
    const res = await uni.request({ url: 'http://localhost:8080/admin/logs', method: 'GET' });
    if (res.statusCode === 200) {
      const clsId = nameToId.value[currentClass.value];
      console.log('Raw logs response:', res.data.logs);
      console.log('Current class:', currentClass.value);
      console.log('Class ID mapping:', nameToId.value);
      console.log('Filtered class ID:', clsId);
      
      const items = (res.data.logs || [])
        .filter(l => l.classId === clsId)
        .map(l => {
          console.log('Processing log entry:', l);
          const ts = (l.CreatedAt || '').replace('T',' ').split('.')[0];
          const operator = l.operator || l.Operator || '';
          const message = l.message || l.Message || '';
          console.log('Timestamp:', ts, 'Operator:', operator, 'Message:', message);
          if (operator) {
            return `${ts} - ${operator}: ${message}`;
          } else {
            return `${ts} - ${message}`;
          }
        })
        .sort((a, b) => {
          // 按时间戳排序，最新的在顶端
          const timeA = new Date(a.split(' - ')[0]);
          const timeB = new Date(b.split(' - ')[0]);
          return timeB - timeA;
        });
      console.log('Final log items:', items);
      logs.value = items;
    }
  } catch (e) {
    console.error('Error loading logs:', e);
  }
};

// 根据时间段位置获取课程名称
const getCourseName = (row, col) => {
  const match = scheduleData.value.find(schedule => 
    schedule.timeSlotRow === row && 
    schedule.timeSlotCol === col
  );
  
  return match ? match.course.name : '';
};

// 获取星期名称
const getDayName = (col) => {
  const days = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
  return days[col] || 'Unknown';
};

// 获取时间段名称
const getPeriodName = (row) => {
  const periods = ['Morning Slot 1', 'Morning Slot 2', 'Afternoon Slot 1', 'Afternoon Slot 2', 'Evening Slot'];
  return periods[row] || 'Unknown';
};

// period行标题：教师显示具体时间
const getPeriodTime = (row) => {
  const times = [
    '08:00 - 08:45',
    '09:00 - 09:45',
    '14:00 - 14:45',
    '15:00 - 15:45',
    '19:00 - 19:45'
  ];
  return times[row] || '';
};

// 处理课程槽点击
const handleSlotClick = (row, col) => {
  // 切换当前slot的时间提示
  if (tooltipSlot.value && tooltipSlot.value.row === row && tooltipSlot.value.col === col) {
    tooltipSlot.value = null;
  } else {
    tooltipSlot.value = { row, col };
  }
  if (isSelectingTarget.value) {
    // 正在选择目标位置
    if (isSlotBusy(row, col)) {
      uni.showToast({ title: 'Teacher busy at this time', icon: 'none' });
      return;
    }
    moveCourseToTarget(row, col);
  } else if (selectedCourse.value) {
    // 如果已经有选中的课程，点击空白地方应该取消选择
    const courseName = getCourseName(row, col);
    if (!courseName) {
      // 点击空白地方，取消当前选择
      cancelSelection();
    }
  } else {
    // 没有选中课程时，选择源课程或添加新课程
    const courseName = getCourseName(row, col);
    if (courseName) {
      selectedCourse.value = { 
        row, 
        col, 
        courseName,
        sourceWeek: currentWeek.value // 保存源周数
      };
    } else {
      // 空时间段，显示添加新课程选项
      if (isSlotBusy(row, col)) {
        uni.showToast({ title: 'Teacher busy at this time', icon: 'none' });
        return;
      }
      newCourseSlot.value = { row, col, week: currentWeek.value };
      isAddingNewCourse.value = true;
    }
  }
};

// 根据行列返回固定时间文本
const getSlotTime = (row, col) => {
  const dayNames = ['Monday','Tuesday','Wednesday','Thursday','Friday','Saturday','Sunday'];
  const periods = [
    '08:00 - 08:45', // Morning Slot 1
    '09:00 - 09:45', // Morning Slot 2
    '14:00 - 14:45', // Afternoon Slot 1
    '15:00 - 15:45', // Afternoon Slot 2
    '19:00 - 19:45'  // Evening Slot
  ];
  const d = dayNames[col] || 'Unknown';
  const t = periods[row] || 'Unknown';
  return `${d} · ${t}`;
};

// 保留 getSlotTime 供提示气泡使用

// 取消选择
const cancelSelection = () => {
  selectedCourse.value = null;
  isSelectingTarget.value = false;
};

// 取消添加新课程
const cancelAddNewCourse = () => {
  isAddingNewCourse.value = false;
  newCourseSlot.value = null;
  newCourseName.value = ''; // 清空输入框
};

// 添加新课程
const addNewCourse = async (courseName) => {
  if (!courseName.trim()) {
    uni.showToast({ title: 'Please enter a course name', icon: 'none' });
    return;
  }
  
  if (!newCourseSlot.value) return;
  
  const userInfo = uni.getStorageSync('userInfo');
  const operator = userInfo ? userInfo.username : 'unknown';
  console.log('Adding course with operator:', operator);
  console.log('User info from storage:', userInfo);
  
  try {
    const requestData = {
      className: currentClass.value,
      courseName: courseName.trim(),
      weekNumber: newCourseSlot.value.week,
      timeSlotRow: newCourseSlot.value.row,
      timeSlotCol: newCourseSlot.value.col,
      operator: operator
    };
    console.log('Request data for add course:', requestData);
    
    const response = await uni.request({
      url: 'http://localhost:8080/api/schedule/add',
      method: 'POST',
      header: {
        'Content-Type': 'application/json'
      },
      data: requestData
    });

    if (response.statusCode === 200) {
      uni.showToast({ title: 'Course added successfully', icon: 'success' });
      isAddingNewCourse.value = false;
      newCourseSlot.value = null;
      loadSchedule(); // 重新加载课程表
      loadLogs(); // 重新加载日志
    } else {
      uni.showToast({ title: response.data.error || 'Failed to add course', icon: 'none' });
    }
  } catch (error) {
    console.error('Failed to add course:', error);
    uni.showToast({ title: 'Failed to add course', icon: 'none' });
  }
};

// 删除课程
const deleteCourse = async () => {
  if (!selectedCourse.value) return;
  
  const userInfo = uni.getStorageSync('userInfo');
  const operator = userInfo ? userInfo.username : 'unknown';
  console.log('Deleting course with operator:', operator);
  
  try {
    const requestData = {
      className: currentClass.value,
      weekNumber: selectedCourse.value.sourceWeek, // 删除源周数的课程
      timeSlotRow: selectedCourse.value.row,
      timeSlotCol: selectedCourse.value.col,
      operator: operator
    };
    console.log('Request data for delete course:', requestData);
    
    const response = await uni.request({
      url: 'http://localhost:8080/api/schedule/delete',
      method: 'DELETE',
      header: {
        'Content-Type': 'application/json'
      },
      data: requestData
    });

    if (response.statusCode === 200) {
      uni.showToast({ title: 'Course deleted successfully', icon: 'success' });
      selectedCourse.value = null;
      loadSchedule(); // 重新加载课程表
      loadLogs(); // 重新加载日志
    } else {
      uni.showToast({ title: response.data.error || 'Failed to delete course', icon: 'none' });
    }
  } catch (error) {
    console.error('Failed to delete course:', error);
    uni.showToast({ title: 'Failed to delete course', icon: 'none' });
  }
};

// 开始选择目标位置
const startTargetSelection = () => {
  isSelectingTarget.value = true;
};

// 取消目标选择
const cancelTargetSelection = () => {
  isSelectingTarget.value = false;
};

// 移动课程到目标位置
const moveCourseToTarget = async (targetRow, targetCol) => {
  if (!selectedCourse.value) return;
  
  // 检查目标位置是否已有课程
  const targetCourseName = getCourseName(targetRow, targetCol);
  if (targetCourseName) {
    uni.showToast({ title: 'Target slot is not empty', icon: 'none' });
    return;
  }
  
  // 直接使用当前周数作为目标周数
  const targetWeek = currentWeek.value;
  
  const userInfo = uni.getStorageSync('userInfo');
  const operator = userInfo ? userInfo.username : 'unknown';
  console.log('Moving course with operator:', operator);
  
  try {
    const requestData = {
      className: currentClass.value,
      sourceWeek: selectedCourse.value.sourceWeek, // 使用选中的源周数
      sourceRow: selectedCourse.value.row,
      sourceCol: selectedCourse.value.col,
      targetWeek: targetWeek,
      targetRow: targetRow,
      targetCol: targetCol,
      operator: operator
    };
    console.log('Request data for move course:', requestData);
    
    const response = await uni.request({
      url: 'http://localhost:8080/api/schedule/move',
      method: 'POST',
      header: {
        'Content-Type': 'application/json'
      },
      data: requestData
    });

    if (response.statusCode === 200) {
      uni.showToast({ title: 'Course moved successfully', icon: 'success' });
      selectedCourse.value = null;
      isSelectingTarget.value = false;
      loadSchedule(); // 重新加载课程表
      loadLogs(); // 重新加载日志
      loadBusySlots(); // 刷新教师忙碌占位（若为教师）
    } else {
      uni.showToast({ title: response.data.error || 'Failed to move course', icon: 'none' });
    }
  } catch (error) {
    console.error('Failed to move course:', error);
    uni.showToast({ title: 'Failed to move course', icon: 'none' });
  }
};



const toggleSidebar = () => {
  sidebarExpanded.value = !sidebarExpanded.value;
};

const closeSidebar = () => {
  sidebarExpanded.value = false;
};

// 显示teacher占位符弹窗
const openTeacherPlanModal = () => {
  planClassName.value = '';
  planWeekType.value = 'continuous';
  planStartWeek.value = currentWeek.value;
  planEndWeek.value = currentWeek.value;
  planSelectedWeeks.value = [];
  planSelectedSlots.value = new Set();
  showTeacherPlanModal.value = true;
};

const closeTeacherPlanModal = () => {
  showTeacherPlanModal.value = false;
};

const togglePlanWeek = (w) => {
  const idx = planSelectedWeeks.value.indexOf(w);
  if (idx >= 0) planSelectedWeeks.value.splice(idx, 1);
  else planSelectedWeeks.value.push(w);
};

const isPlanSelected = (row, col) => {
  return planSelectedSlots.value.has(`${row}-${col}`);
};

const togglePlanSlot = (row, col) => {
  const key = `${row}-${col}`;
  const next = new Set(planSelectedSlots.value);
  if (next.has(key)) next.delete(key);
  else next.add(key);
  planSelectedSlots.value = next; // 触发响应
};

const buildWeeks = () => {
  if (planWeekType.value === 'continuous') {
    const res = [];
    for (let w = planStartWeek.value; w <= planEndWeek.value; w++) res.push(w);
    return res;
  }
  return [...planSelectedWeeks.value];
};

const saveTeacherPlan = async () => {
  if (userType.value !== 'teacher') {
    uni.showToast({ title: 'Only teachers can use this', icon: 'none' });
    return;
  }
  if (!planClassName.value.trim()) {
    uni.showToast({ title: 'Please enter class ID', icon: 'none' });
    return;
  }
  const weeks = buildWeeks();
  if (weeks.length === 0) {
    uni.showToast({ title: 'Please select weeks', icon: 'none' });
    return;
  }
  const slots = Array.from(planSelectedSlots.value).map(k => {
    const [r, c] = k.split('-').map(n => parseInt(n));
    return { row: r, col: c };
  });
  if (slots.length === 0) {
    uni.showToast({ title: 'Please select at least one slot', icon: 'none' });
    return;
  }
  try {
    const res = await uni.request({
      url: 'http://localhost:8080/api/teacher/teach-schedule/overwrite-batch',
      method: 'POST',
      header: { 'Content-Type': 'application/json' },
      data: {
        teacherId: userID.value,
        className: planClassName.value.trim(),
        weeks,
        slots,
        operator: uni.getStorageSync('userInfo').username // 添加操作人
      }
    });
    if (res.statusCode === 200) {
      uni.showToast({ title: 'Saved', icon: 'success' });
      showTeacherPlanModal.value = false;
      await loadBusySlots();
    } else {
      uni.showToast({ title: res.data.error || 'Save failed', icon: 'none' });
    }
  } catch (e) {
    uni.showToast({ title: 'Save failed', icon: 'none' });
  }
};

const clearAllTeacherPlans = async () => {
  if (userType.value !== 'teacher') {
    uni.showToast({ title: 'Only teachers can use this', icon: 'none' });
    return;
  }
  const wasOpen = showTeacherPlanModal.value;
  // 暂时关闭录入弹窗，避免遮挡系统确认框
  showTeacherPlanModal.value = false;
  uni.showModal({
    title: 'Confirm',
    content: 'This will clear ALL your plans across all classes and weeks. Continue?',
    success: async (res) => {
      if (!res.confirm) {
        if (wasOpen) showTeacherPlanModal.value = true;
        return;
      }
      try {
        const result = await uni.request({
          url: `http://localhost:8080/api/teacher/teach-schedule/all?teacherId=${encodeURIComponent(userID.value)}`,
          method: 'DELETE',
        });
        if (result.statusCode === 200) {
          uni.showToast({ title: 'All plans cleared', icon: 'success' });
          // 保持关闭状态
          await loadBusySlots();
        } else {
          uni.showToast({ title: result.data.error || 'Clear failed', icon: 'none' });
          if (wasOpen) showTeacherPlanModal.value = true;
        }
      } catch (e) {
        uni.showToast({ title: 'Clear failed', icon: 'none' });
        if (wasOpen) showTeacherPlanModal.value = true;
      }
    },
    fail: () => {
      if (wasOpen) showTeacherPlanModal.value = true;
    }
  });
};

const isSlotBusy = (row, col) => {
  return busySlots.value.has(`${row}-${col}`);
};

// 查看计划（只读）
const viewerWeek = ref(1);
const viewerBusy = ref({}); // key -> [className,...]
const openTeacherPlanViewer = async () => {
  viewerWeek.value = currentWeek.value;
  await loadViewerBusy();
  showTeacherPlanViewer.value = true;
};
const closeTeacherPlanViewer = () => { showTeacherPlanViewer.value = false; };
const prevViewerWeek = async () => { if (viewerWeek.value > 1) viewerWeek.value--; await loadViewerBusy(); };
const nextViewerWeek = async () => { if (viewerWeek.value < 20) viewerWeek.value++; await loadViewerBusy(); };
const isViewerBusy = (row, col) => {
  const key = `${row}-${col}`;
  const arr = viewerBusy.value[key];
  return Array.isArray(arr) && arr.length > 0;
};
const getViewerClassText = (row, col) => {
  const key = `${row}-${col}`;
  const arr = viewerBusy.value[key] || [];
  // 将每个 className 格式化为 “第3位和第4位/最后一位”，不足位数的做保护
  const fmt = (nameStr) => {
    const s = String(nameStr);
    const third = s.length >= 3 ? s[2] : '';
    const fourth = s.length >= 4 ? s[3] : '';
    const last = s.length >= 1 ? s[s.length - 1] : '';
    return `${third}${fourth}/${last}`;
  };
  return arr.map(fmt).join('/');
};
const loadViewerBusy = async () => {
  try {
    const res = await uni.request({
      url: `http://localhost:8080/api/teacher/teach-busy?teacherId=${encodeURIComponent(userID.value)}&week=${viewerWeek.value}`,
      method: 'GET'
    });
    if (res.statusCode === 200) {
      const map = {};
      (res.data.slots || []).forEach(s => {
        const key = `${s.row}-${s.col}`;
        if (!map[key]) map[key] = [];
        // 使用 className 作为来源字符串
        const text = (s.className !== undefined && s.className !== null) ? String(s.className) : '';
        if (text && !map[key].includes(text)) map[key].push(text);
      });
      viewerBusy.value = map;
    }
  } catch (e) {}
};

const loadBusySlots = async () => {
  if (userType.value !== 'teacher' || !userID.value) return;
  try {
    const res = await uni.request({
      url: `http://localhost:8080/api/teacher/teach-busy?teacherId=${encodeURIComponent(userID.value)}&week=${currentWeek.value}&excludeClassName=${encodeURIComponent(currentClass.value || '')}`,
      method: 'GET'
    });
    if (res.statusCode === 200) {
      const set = new Set();
      (res.data.slots || []).forEach(s => set.add(`${s.row}-${s.col}`));
      busySlots.value = set;
    }
  } catch (e) {}
};

// 获取用户类型
const getUserType = () => {
  const userInfo = uni.getStorageSync('userInfo');
  console.log('Getting user type from storage:', userInfo);
  if (userInfo && userInfo.userType) {
    userType.value = userInfo.userType;
    userID.value = userInfo.userID || '';
    console.log('Set user type to:', userType.value);
    console.log('Set user ID to:', userID.value);
  } else {
    // 如果没有存储的用户信息，默认为user类型
    userType.value = 'user';
    console.log('No user info found, defaulting to user type');
  }
};

onMounted(() => {
  getUserType();
  loadClasses();
  loadLogs();
  loadBusySlots();
});

// 页面显示时刷新数据
onLoad(() => {
  getUserType();
  loadClasses();
  loadBusySlots();
});
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Patrick+Hand&display=swap');

.main-container {
  display: flex;
  height: 100vh;
  background: #fdfdfb;
  font-family: 'Patrick Hand', cursive;
  padding: 16px;
  box-sizing: border-box;
  background-image: url('https://www.transparenttextures.com/patterns/paper-fibers.png');
  background-size: auto;
  position: relative;
  overflow: hidden; /* 防止侧边栏溢出 */
}

.sidebar-toggle {
  position: fixed;
  top: 20px;
  left: 20px;
  z-index: 1000;
  cursor: pointer;
  padding: 12px;
  background: #fff;
  border: 2px solid #333;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  display: flex; /* 使用flex布局 */
  align-items: center; /* 垂直居中 */
  gap: 8px; /* 汉堡菜单和文字之间的间距 */
  transition: all 0.3s ease;
}

.sidebar-toggle:hover {
  background: #f8f9fa;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
  transform: scale(1.1);
  border-color: #007bff;
}

.hamburger-icon {
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  width: 20px;
  height: 16px;
}

.hamburger-icon span {
  display: block;
  width: 100%;
  height: 3px;
  background: #333;
  border-radius: 2px;
  transition: all 0.3s ease;
}

/* 汉堡菜单动画效果 */
.sidebar-toggle:hover .hamburger-icon span:nth-child(1) {
  transform: rotate(45deg) translate(5px, 5px);
}

.sidebar-toggle:hover .hamburger-icon span:nth-child(2) {
  opacity: 0;
}

.sidebar-toggle:hover .hamburger-icon span:nth-child(3) {
  transform: rotate(-45deg) translate(7px, -6px);
}

.toggle-text {
  font-size: 14px;
  letter-spacing: 0.5px;
  color: #333;
  font-weight: 500;
}

/* 背景遮罩 */
.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  z-index: 100;
  opacity: 0;
  visibility: hidden;
  transition: all 0.3s ease;
}

.sidebar-overlay.active {
  opacity: 1;
  visibility: visible;
}

.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: 250px;
  height: 100vh;
  border-right: 2px dashed #000;
  padding: 20px;
  padding-top: 80px; /* 增加顶部内边距，避免被汉堡菜单按钮遮挡 */
  display: flex;
  flex-direction: column;
  background: #fdfdfb;
  z-index: 200;
  transform: translateX(-100%);
  transition: transform 0.3s ease-in-out;
  box-shadow: 2px 0 10px rgba(0, 0, 0, 0.1);
}

.sidebar-expanded {
  transform: translateX(0);
}

.sidebar h2 {
  margin-bottom: 10px;
}

.sidebar ul {
  list-style: none;
  padding: 0;
  flex: 1;
  overflow-y: auto;
  padding-bottom: 80px; /* 预留底部空间，避免被底部按钮遮挡 */
}

.sidebar li {
  cursor: pointer;
  padding: 5px;
  border: 1px solid #000;
  margin-bottom: 10px;
  transition: background 0.2s;
}

.sidebar li:hover {
  background: rgba(0, 0, 0, 0.05);
}

/* 移除侧边栏底部按钮样式（不再使用） */

/* Teacher专用占位符样式 */
.teacher-placeholder {
  background: #f8f9fa;
  border: 2px dashed #007bff;
  border-radius: 8px;
  padding: 15px;
  margin-bottom: 15px;
  text-align: center;
}

.placeholder-content h4 {
  margin: 0 0 8px 0;
  color: #007bff;
  font-size: 1em;
}

.placeholder-content p {
  margin: 0 0 10px 0;
  color: #666;
  font-size: 0.9em;
}

.placeholder-icon {
  font-size: 2em;
  margin-top: 10px;
}

/* Teacher弹窗样式 */
.teacher-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  display: flex;
  justify-content: center;
  align-items: center;
}

.teacher-modal-content {
  background: #fff;
  border-radius: 12px;
  padding: 30px;
  max-width: 400px;
  width: 90%;
  text-align: center;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  max-height: 80vh; /* 弹窗内容最大高度 */
  overflow-y: auto;  /* 竖向滚动 */
  -webkit-overflow-scrolling: touch;
}

.close-btn {
  margin-top: 20px;
  padding: 10px 20px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.3s;
}

.close-btn:hover {
  background: #0056b3;
}

.schedule-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 20px;
  min-width: 0;
  width: 100%;
  height: 100vh;
  box-sizing: border-box;
  transition: all 0.3s ease;
  max-width: 100vw;
  padding: 0;
}

/* 响应式设计：在小屏幕上调整布局 */
@media (max-width: 768px) {
  .main-container {
    padding: 4px;
    height: 100vh;
    overflow-y: auto;
    overflow-x: hidden;
  }
  
  .sidebar {
    width: 280px;
  }
  
  .schedule-panel {
    padding: 4px;
    height: auto;
    min-height: 100vh;
    overflow: visible;
    display: flex;
    flex-direction: column;
  }
  
  .schedule-top {
    flex-direction: column;
    gap: 4px;
    margin-bottom: 8px;
    flex-shrink: 0;
  }
  
  .schedule-top .title h2 {
    font-size: 1em;
    margin: 0;
    line-height: 1.2;
  }
  
  .schedule-body {
    flex: 1;
    min-height: 0;
    padding: 4px;
    overflow: hidden;
    border: 1px dashed #000;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .schedule-table {
    width: 90%;
    max-width: 400px;
    font-size: 0.8em;
    table-layout: fixed;
  }
  
  .schedule-table th,
  .schedule-table td {
    width: auto;
    min-width: 20px;
    font-size: 0.8em;
    padding: 2px;
  }
  
  .schedule-table th {
    height: 30px;
    min-height: 30px;
    padding: 2px;
  }
  
  .schedule-table td {
    height: 30px;
    min-height: 30px;
  }
  
  .schedule-table td:first-child {
    width: 50px;
    min-width: 50px;
    padding: 2px;
    font-size: 0.7em;
  }
  
  .slot-btn {
    width: 25px !important;
    height: 30px !important;
    min-width: 25px !important;
    min-height: 30px !important;
    padding: 2px;
    font-size: 0.7em;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .course-swap-panel,
  .target-selection-panel,
  .add-course-panel {
    padding: 5px;
    margin-bottom: 5px;
    flex-shrink: 0;
  }
  
  .course-swap-panel h3,
  .target-selection-panel h3,
  .add-course-panel h3 {
    font-size: 0.9em;
    margin-bottom: 5px;
  }
  
  .swap-info p,
  .target-info p,
  .add-course-info p {
    font-size: 0.8em;
    margin-bottom: 5px;
  }
  
  .swap-actions,
  .add-course-actions {
    flex-direction: column;
    gap: 5px;
  }
  
  .swap-actions button,
  .add-course-actions button {
    padding: 8px 12px;
    font-size: 12px;
  }
  
  .add-course-form {
    flex-direction: column;
    gap: 5px;
  }
  
  .course-name-input {
    padding: 8px;
    font-size: 16px;
  }
  
  .log-panel {
    max-height: 80px;
    padding: 5px;
    flex-shrink: 0;
    overflow: auto;
  }
  
  .log-panel h3 {
    font-size: 0.9em;
    margin-bottom: 5px;
  }
  
  .log-content {
    font-size: 0.7em;
  }
  
  .teacher-placeholder {
    padding: 10px;
    margin-bottom: 10px;
  }
  
  .placeholder-content h4 {
    font-size: 0.9em;
  }
  
  .placeholder-content p {
    font-size: 0.8em;
  }
  
  .placeholder-icon {
    font-size: 1.5em;
  }
  
  .teacher-modal-content {
    padding: 20px;
    max-width: 350px;
  }
  
  .close-btn {
    margin-top: 15px;
    padding: 8px 16px;
    font-size: 13px;
  }
}

/* 超小屏幕优化 */
@media (max-width: 480px) {
  .main-container {
    padding: 2px;
  }
  
  .schedule-panel {
    padding: 2px;
  }
  
  .schedule-top {
    gap: 2px;
    margin-bottom: 4px;
  }
  
  .schedule-top .title h2 {
    font-size: 0.9em;
  }
  
  .schedule-body {
    padding: 2px 0;
  }
  
  .schedule-table {
    width: 100%;
    font-size: 0.7em;
  }
  
  .schedule-table th,
  .schedule-table td {
    width: auto;
    min-width: 20px;
    font-size: 0.7em;
    padding: 2px;
  }
  
  .schedule-table th {
    height: 25px;
    min-height: 25px;
    padding: 2px;
  }
  
  .schedule-table td {
    height: 25px;
    min-height: 25px;
  }
  
  .schedule-table td:first-child {
    width: 40px;
    min-width: 40px;
    padding: 2px;
    font-size: 0.6em;
  }
  
  .slot-btn {
    width: 20px !important;
    height: 25px !important;
    min-width: 20px !important;
    min-height: 25px !important;
    padding: 2px;
    font-size: 0.6em;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .course-swap-panel,
  .target-selection-panel,
  .add-course-panel {
    padding: 3px;
    margin-bottom: 3px;
  }
  
  .course-swap-panel h3,
  .target-selection-panel h3,
  .add-course-panel h3 {
    font-size: 0.8em;
    margin-bottom: 3px;
  }
  
  .swap-info p,
  .target-info p,
  .add-course-info p {
    font-size: 0.7em;
    margin-bottom: 3px;
  }
  
  .swap-actions button,
  .add-course-actions button {
    padding: 6px 10px;
    font-size: 11px;
  }
  
  .log-panel {
    max-height: 60px;
    padding: 3px;
  }
  
  .log-panel h3 {
    font-size: 0.8em;
    margin-bottom: 3px;
  }
  
  .log-content {
    font-size: 0.6em;
  }
  
  .teacher-placeholder {
    padding: 8px;
    margin-bottom: 8px;
  }
  
  .placeholder-content h4 {
    font-size: 0.8em;
  }
  
  .placeholder-content p {
    font-size: 0.7em;
  }
  
  .placeholder-icon {
    font-size: 1.2em;
  }
  
  .teacher-modal-content {
    padding: 15px;
    max-width: 300px;
  }
  
  .close-btn {
    margin-top: 12px;
    padding: 6px 12px;
    font-size: 12px;
  }
}

.class-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.add-btn {
  background: none;
  border: none;
  font-size: 1.2em;
  width: 24px;
  height: 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s;
}
.add-btn:hover {
  transform: scale(1.2);
}

.schedule-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.schedule-top .title {
  flex: 1;
  text-align: center;
}

.schedule-top .controls button {
  background: none;
  border: none;
  font-size: 1.2em;
  cursor: pointer;
  padding: 2px 6px;
  transition: transform 0.2s;
  line-height: 1;
  height: auto;
}

.schedule-top .controls button:hover {
  transform: scale(1.2);
}

.schedule-body {
  flex: 1;
  border: 2px dashed #000;
  margin-bottom: 10px;
  padding: 15px 0;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.course-swap-panel {
  border-top: 2px dashed #000;
  padding: 10px;
  background: #f8f9fa;
  margin-bottom: 10px;
}

.course-swap-panel h3 {
  margin-bottom: 10px;
  color: #333;
}

.swap-info p {
  margin-bottom: 10px;
  font-weight: bold;
}

.swap-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.swap-actions button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.cancel-btn {
  background: #6c757d;
  color: white;
}

.cancel-btn:hover {
  background: #5a6268;
}

.delete-btn {
  background: #dc3545;
  color: white;
}

.delete-btn:hover {
  background: #c82333;
}

.swap-btn {
  background: #007bff;
  color: white;
}

.swap-btn:hover {
  background: #0056b3;
}

.target-selection-panel {
  border-top: 2px dashed #000;
  padding: 10px;
  background: #fff3cd;
  margin-bottom: 10px;
}

.target-selection-panel h3 {
  margin-bottom: 10px;
  color: #856404;
}

.target-info p {
  margin-bottom: 10px;
  color: #856404;
}

.add-course-panel {
  border-top: 2px dashed #000;
  padding: 10px;
  background: #e0f2f7;
  margin-bottom: 10px;
}

.add-course-panel h3 {
  margin-bottom: 10px;
  color: #007bff;
}

.add-course-info p {
  margin-bottom: 10px;
  font-weight: bold;
  color: #333;
}

.add-course-form {
  display: flex;
  gap: 10px;
  align-items: center;
}

.course-name-input {
  flex: 1;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-family: 'Patrick Hand', cursive;
  font-size: 0.9em;
}

.add-course-actions {
  display: flex;
  gap: 10px;
}

.add-course-actions button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
}

.confirm-add-btn {
  background: #28a745;
  color: white;
}

.confirm-add-btn:hover {
  background: #218838;
}

.log-panel {
  border-top: 2px dashed #000;
  padding: 10px;
  max-height: 150px;
  overflow-y: auto;
}

.log-panel h3 {
  margin-bottom: 10px;
}

.log-content {
  font-size: 0.9em;
}

.schedule-table {
  width: 100%;
  max-width: 100vw;
  margin: 0 auto;
  border-collapse: collapse;
  font-size: 1em;
  table-layout: fixed;
  text-align: center;
}

.schedule-table th,
.schedule-table td {
  border: 1px solid #000;
  padding: 0;
  font-size: 0.8em;
  width: 40px;
  min-width: 40px;
  height: 50px;
  min-height: 50px;
}

.schedule-table thead {
  background-color: #f5f5f5;
}

.schedule-table th {
  background-color: #eaeaea;
  font-weight: bold;
  font-size: 0.8em;
  height: 60px;
  min-height: 60px;
  padding: 8px;
}

.schedule-table td {
  height: 60px;
  min-height: 60px;
  vertical-align: middle;
  padding: 0;
}

.schedule-table td:first-child {
  font-weight: bold;
  background-color: #f8f8f8;
  width: 60px;
  min-width: 60px;
  padding: 0;
  font-size: 0.95em;
}

.slot-btn {
  width: 60px !important;
  min-width: 60px !important;
  height: 60px !important;
  min-height: 60px !important;
  font-size: 0.75em;
  line-height: 1.2;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  word-wrap: break-word;
  word-break: break-word;
  overflow: hidden;
  text-align: center;
}

.slot-btn.selected {
  background: #4CAF50;
  color: white;
  border-color: #4CAF50;
}

.slot-btn.target {
  background: #FF9800;
  color: white;
  border-color: #FF9800;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); }
}

.slot-btn {
  font-size: 0.75em;
  line-height: 1.2;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  word-wrap: break-word;
  word-break: break-word;
  overflow: hidden;
  text-align: center;
}
.slot-btn:hover {
  background-color: #f0f0f0;
}

.slot-btn.busy {
  background: rgba(220, 53, 69, 0.2);
  border-color: #dc3545;
}

/* 槽位时间提示样式 */
.slot-tooltip {
  position: absolute;
  margin-top: 4px;
  left: 50%;
  transform: translateX(-50%);
  background: #333;
  color: #fff;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  white-space: nowrap;
  z-index: 10;
  box-shadow: 0 2px 6px rgba(0,0,0,0.2);
}

td { position: relative; }

/* 弹窗内计划表样式复用 register_class 弹窗风格 */
.teacher-plan-header h4 {
  margin: 0 0 8px 0;
  color: #007bff;
}
.teacher-plan-body .input-group label { font-weight: bold; }
.plan-input { width: 100%; padding: 8px; border: 1px solid #ccc; border-radius: 4px; }
.week-type-buttons { display: flex; gap: 8px; margin: 8px 0; }
.week-type-btn { flex: 1; padding: 8px; border: 1px solid #ddd; background: #fff; border-radius: 6px; }
.week-type-btn.active { border-color: #007bff; background: #e3f2fd; color: #007bff; }
.week-range { display: flex; gap: 10px; }
.range-input { flex: 1; }
.week-select { width: 100%; padding: 6px; }
.week-grid { display: grid; grid-template-columns: repeat(5, 1fr); gap: 6px; max-width: 260px; }
.week-number-btn { aspect-ratio: 1; border: 1px solid #ddd; border-radius: 6px; }
.week-number-btn.selected { background: #007bff; color: #fff; border-color: #007bff; }
.plan-grid-wrapper { margin-top: 10px; }
.plan-grid-title { font-weight: bold; margin-bottom: 6px; }
.plan-table { width: 100%; border-collapse: collapse; table-layout: fixed; }
.plan-table th, .plan-table td { border: 1px solid #000; text-align: center; }
.plan-slot { width: 40px !important; height: 50px !important; min-width: 40px !important; min-height: 50px !important; }
.teacher-plan-footer { margin-top: 10px; display: flex; gap: 10px; justify-content: flex-end; }
.save-btn { background: #28a745; color: #fff; border: none; padding: 10px 16px; border-radius: 6px; }
.teacher-plan-footer .clear-btn { background: #dc3545; color: #fff; border: none; padding: 10px 16px; border-radius: 6px; }
</style>