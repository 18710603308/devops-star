import { createApp } from "vue";
import ElementPlus from "element-plus";
import "element-plus/dist/index.css";
import "element-plus/theme-chalk/dark/css-vars.css";
import zhCn from "element-plus/dist/locale/zh-cn.mjs";
import App from "./App.vue";
import router from "./router";
import { createPinia } from "pinia";
import * as echarts from "echarts";

const app = createApp(App);

app.use(ElementPlus, { locale: zhCn });
app.use(router);
app.use(createPinia());

// 全局挂载 ECharts
app.config.globalProperties.$echarts = echarts;

app.mount("#app");
