<template>
  <el-config-provider :locale="zhCn" :button-auto-insert-space="false">
    <el-container class="app-container">
      <!-- 桌面端侧边栏 -->
      <el-aside v-if="!isMobile" :width="isCollapse ? '64px' : '220px'" class="app-sidebar">
        <div class="sidebar-logo">
          <span v-if="!isCollapse" class="logo-text">DevOpsStar</span>
          <span v-else class="logo-text">D</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          :collapse="isCollapse"
          router
          background-color="#1d1e2c"
          text-color="#bfcbd9"
          active-text-color="#409EFF"
        >
          <el-menu-item index="/dashboard">
            <el-icon><DataLine /></el-icon>
            <template #title>仪表盘</template>
          </el-menu-item>
          <el-menu-item index="/projects">
            <el-icon><Folder /></el-icon>
            <template #title>项目管理</template>
          </el-menu-item>
          <el-menu-item index="/pipeline">
            <el-icon><Share /></el-icon>
            <template #title>流水线</template>
          </el-menu-item>
          <el-menu-item index="/registry">
            <el-icon><Box /></el-icon>
            <template #title>制品管理</template>
          </el-menu-item>
          <el-menu-item index="/deploy">
            <el-icon><Promotion /></el-icon>
            <template #title>部署管理</template>
          </el-menu-item>
          <el-menu-item index="/monitor">
            <el-icon><Monitor /></el-icon>
            <template #title>监控大屏</template>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <template #title>系统设置</template>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <!-- 移动端抽屉菜单 -->
      <el-drawer
        v-if="isMobile"
        v-model="mobileMenuOpen"
        direction="ltr"
        size="220px"
        :with-header="false"
        class="mobile-drawer"
      >
        <div class="sidebar-logo mobile-logo">
          <span class="logo-text">DevOpsStar</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          router
          background-color="#1d1e2c"
          text-color="#bfcbd9"
          active-text-color="#409EFF"
        >
          <el-menu-item index="/dashboard">
            <el-icon><DataLine /></el-icon>
            <template #title>仪表盘</template>
          </el-menu-item>
          <el-menu-item index="/projects">
            <el-icon><Folder /></el-icon>
            <template #title>项目管理</template>
          </el-menu-item>
          <el-menu-item index="/pipeline">
            <el-icon><Share /></el-icon>
            <template #title>流水线</template>
          </el-menu-item>
          <el-menu-item index="/registry">
            <el-icon><Box /></el-icon>
            <template #title>制品管理</template>
          </el-menu-item>
          <el-menu-item index="/deploy">
            <el-icon><Promotion /></el-icon>
            <template #title>部署管理</template>
          </el-menu-item>
          <el-menu-item index="/monitor">
            <el-icon><Monitor /></el-icon>
            <template #title>监控大屏</template>
          </el-menu-item>
          <el-menu-item index="/settings">
            <el-icon><Setting /></el-icon>
            <template #title>系统设置</template>
          </el-menu-item>
        </el-menu>
      </el-drawer>

      <!-- 主内容区 -->
      <el-container>
        <el-header class="app-header">
          <div class="header-left">
            <!-- 移动端菜单按钮 -->
            <el-icon
              v-if="isMobile"
              class="mobile-menu-btn"
              @click="mobileMenuOpen = true"
            >
              <Menu />
            </el-icon>
            <!-- 桌面端折叠按钮 -->
            <el-icon
              v-else
              class="collapse-btn"
              @click="isCollapse = !isCollapse"
            >
              <Fold v-if="!isCollapse" />
              <Expand v-else />
            </el-icon>
            <el-breadcrumb separator="/">
              <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
              <el-breadcrumb-item>{{ currentRouteTitle }}</el-breadcrumb-item>
            </el-breadcrumb>
          </div>
          <div class="header-right">
            <el-dropdown @command="handleUserCommand">
              <span class="user-info">
                <el-icon><User /></el-icon>
                {{ userStore.username || '管理员' }}
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="settings">系统设置</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        <el-main class="app-main">
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </el-config-provider>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import { useUserStore } from '@/stores/user'
import {
  DataLine, Folder, Share, Box, Promotion, Monitor, Setting,
  User, ArrowDown, Fold, Expand, Menu
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)
const activeMenu = computed(() => route.path)
const mobileMenuOpen = ref(false)
const isMobile = ref(false)

const routeTitleMap: Record<string, string> = {
  '/dashboard': '仪表盘',
  '/projects': '项目管理',
  '/pipeline': '流水线编排',
  '/registry': '制品管理',
  '/deploy': '部署管理',
  '/monitor': '监控大屏',
  '/settings': '系统设置',
}

const currentRouteTitle = computed(() => routeTitleMap[route.path] || '首页')

const handleUserCommand = (command: string) => {
  if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  } else if (command === 'profile') {
    router.push('/settings')
  } else if (command === 'settings') {
    router.push('/settings')
  }
}

// 检测移动端
const checkMobile = () => {
  isMobile.value = window.innerWidth < 768
  if (isMobile.value) {
    isCollapse.value = false
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

// 登录状态校验
if (route.path !== '/login' && !userStore.token) {
  router.push('/login')
}
</script>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
}
.app-container {
  height: 100vh;
}
.app-sidebar {
  background-color: #1d1e2c;
  overflow-y: auto;
}
.sidebar-logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #409EFF;
  font-size: 20px;
  font-weight: bold;
  border-bottom: 1px solid #2d2e3c;
}
.logo-text {
  letter-spacing: 2px;
}
.app-header {
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 56px !important;
  box-shadow: 0 1px 4px rgba(0,0,0,0.08);
}
.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.collapse-btn, .mobile-menu-btn {
  font-size: 20px;
  cursor: pointer;
  color: #606266;
}
.header-right {
  display: flex;
  align-items: center;
}
.user-info {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  color: #606266;
  font-size: 14px;
}
.app-main {
  background: #f0f2f5;
  padding: 20px;
  min-height: 0;
  overflow-y: auto;
}

/* 移动端适配 */
@media screen and (max-width: 767px) {
  .app-header {
    padding: 0 12px;
  }
  .header-left {
    gap: 8px;
  }
  .app-main {
    padding: 12px;
  }
  .mobile-menu-btn {
    font-size: 24px;
  }
}

/* 平板适配 */
@media screen and (min-width: 768px) and (max-width: 1024px) {
  .app-main {
    padding: 16px;
  }
}

/* 大屏适配 */
@media screen and (min-width: 1920px) {
  .app-main {
    max-width: 1600px;
    margin: 0 auto;
  }
}

/* 移动端抽屉菜单样式 */
.mobile-drawer .el-drawer__body {
  padding: 0;
}
.mobile-logo {
  border-bottom: 1px solid #2d2e3c;
}
</style>
