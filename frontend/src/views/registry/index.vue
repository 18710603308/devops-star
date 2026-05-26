<template>
  <div class="registry-page">
    <el-card class="page-header">
      <div class="header-actions">
        <h2>制品管理</h2>
        <el-button type="primary" @click="showPush = true">
          <el-icon><Plus /></el-icon>
          推送镜像
        </el-button>
      </div>
    </el-card>

    <el-card class="table-card">
      <el-table :data="imageList" stripe style="width: 100%">
        <el-table-column prop="name" label="镜像名称" />
        <el-table-column prop="tag" label="标签" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.tag }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="120" />
        <el-table-column prop="push_time" label="推送时间" width="180" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="danger" size="small" @click="deleteImage(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showPush" title="推送镜像" width="520px">
      <el-form label-width="100px">
        <el-form-item label="镜像名称">
          <el-input v-model="pushForm.name" placeholder="请输入镜像名称" />
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="pushForm.tag" placeholder="latest" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showPush = false">取消</el-button>
        <el-button type="primary" @click="handlePush">推送</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";

const imageList = ref([
  { name: "devops-star/frontend", tag: "latest", size: "45.2 MB", push_time: "2026-05-27 00:30:00" },
  { name: "devops-star/backend", tag: "v1.0.0", size: "32.1 MB", push_time: "2026-05-27 00:20:00" },
]);
const showPush = ref(false);
const pushForm = ref({ name: "", tag: "latest" });

const loadImages = async () => {
  try {
    const res = await import("@/api/registry").then(m => m.getImages());
    imageList.value = res.data || imageList.value;
  } catch {}
};

const handlePush = () => {
  ElMessage.success("镜像推送指令已生成，请在本地执行 Docker Push");
  showPush.value = false;
};

const deleteImage = async (row: any) => {
  await ElMessageBox.confirm("确定删除该镜像？", "提示", { type: "warning" });
  try {
    await import("@/api/registry").then(m => m.deleteImage(row.name, row.tag));
    ElMessage.success("删除成功");
    loadImages();
  } catch {}
};

onMounted(loadImages);
</script>

<style scoped>
.registry-page { display: flex; flex-direction: column; gap: 20px; }
.page-header { padding: 20px; }
.header-actions { display: flex; justify-content: space-between; align-items: center; }
.header-actions h2 { margin: 0; }
.table-card { padding: 20px; }
</style>
