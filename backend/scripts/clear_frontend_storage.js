// 清理前端本地存储脚本
// 在浏览器控制台或uni-app环境中运行此脚本

console.log('🧹 开始清理前端本地存储...');

try {
  // 清理班级列表
  if (typeof uni !== 'undefined') {
    uni.removeStorageSync('classList');
    console.log('✅ 班级列表已清理');
    
    // 清理其他可能的存储数据
    uni.removeStorageSync('userInfo');
    uni.removeStorageSync('scheduleData');
    uni.removeStorageSync('currentClass');
    uni.removeStorageSync('currentWeek');
    
    console.log('✅ 前端本地存储清理完成');
  } else if (typeof localStorage !== 'undefined') {
    // 浏览器环境
    localStorage.removeItem('classList');
    localStorage.removeItem('userInfo');
    localStorage.removeItem('scheduleData');
    localStorage.removeItem('currentClass');
    localStorage.removeItem('currentWeek');
    
    console.log('✅ 浏览器本地存储清理完成');
  } else {
    console.log('⚠️  无法识别存储环境');
  }
} catch (error) {
  console.error('❌ 清理前端存储时出错:', error);
}

console.log('�� 提示：请刷新页面以查看清理效果'); 