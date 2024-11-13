<template>
  <div class="app-container">
    <el-form ref="queryForm" :model="queryParams" class="query-form" size="small">
      <el-input
        v-model="queryParams.title"
        placeholder="项目名称"
        clearable
        @keyup.enter="getProjectPage()"
      />
      <el-input
        v-model="queryParams.summary"
        placeholder="项目描述"
        clearable
        @keyup.enter="getProjectPage()"
      />

      <el-select v-model="queryParams.type" placeholder="项目类型" clearable>
        <el-option key="0" label="活跃" value="0" />
        <el-option key="1" label="存档" value="1" />
      </el-select>

      <el-select v-model="queryParams.status" placeholder="公开?" clearable>
        <el-option key="0" label="公开" value="0" />
        <el-option key="1" label="隐藏" value="1" />
      </el-select>
      <el-button type="primary" :icon="Search" @click="getProjectPage()">搜索</el-button>
    </el-form>

    <el-row :gutter="10">
      <el-col :span="1.5">
        <el-button type="primary" plain :icon="Plus" @click="handleAdd">新增</el-button>
      </el-col>
    </el-row>

    <el-table v-loading="loading" :data="projectList">
      <el-table-column label="ID" align="center" width="50" prop="id" />
      <el-table-column label="名称" align="center" prop="title" />
      <el-table-column
        label="描述"
        align="center"
        show-overflow-tooltip
        width="200"
        prop="summary"
      />

      <el-table-column label="logo" align="center" width="100" prop="thumbnail" type="img">
        <template v-slot="scope">
          <el-image style="width: 88px" :src="scope.row.thumbnail" fit="cover" />
        </template>
      </el-table-column>

      <el-table-column prop="type" label="类型" align="center">
        <template v-slot="scope">
          <el-tag v-if="scope.row.type === '0'" type="success">活跃</el-tag>
          <el-tag v-else type="info">存档</el-tag>
        </template>
      </el-table-column>

      <el-table-column label="地址" width="300" align="center">
        <template v-slot="scope">
          <el-link :href="scope.row.srcUrl" target="_blank" style="word-break: break-all">
            <span>源地址：</span>
            {{ scope.row.srcUrl }}
          </el-link>
          <br />
          <el-link :href="scope.row.displayUrl" target="_blank" style="word-break: break-all">
            <span>展示地址：</span>
            {{ scope.row.displayUrl }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column prop="beginDate" label="开始日期" width="100" align="center" />

      <el-table-column prop="status" label="公开" align="center">
        <template v-slot="scope">
          <el-switch
            v-model="scope.row.status"
            active-value="0"
            inactive-value="1"
            :loading="statusLoading"
            @change="handleChangeStatus(scope.row)"
          />
        </template>
      </el-table-column>

      <el-table-column prop="isTop" label="置顶" align="center">
        <template v-slot="scope">
          <el-switch
            v-model="scope.row.isTop"
            active-value="1"
            inactive-value="0"
            :loading="statusLoading"
            @change="handleChangeStatus(scope.row)"
          />
        </template>
      </el-table-column>

      <el-table-column prop="orderNum" label="排序" align="center" />

      <el-table-column label="操作" align="center" fixed="right" width="170">
        <template v-slot="scope">
          <el-button type="text" :icon="Edit" @click="handleUpdate(scope.row)" size="small">
            修改
          </el-button>
          <el-button type="text" :icon="Delete" @click="handleDelete(scope.row)" size="small">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:page-size="queryParams.pageSize"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total"
      :page-sizes="[5, 10, 20, 30, 40]"
      v-model:current-page="queryParams.pageNum"
      @current-change="getProjectPage"
      @size-change="getProjectPage()"
    />

    <el-dialog :title="title" v-model="open" class="dialog-form">
      <el-form ref="form" :rules="rules" label-width="auto" label-position="top">
        <el-form-item label="名称" prop="name">
          <el-input v-model="project.title" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="project.summary" type="textarea" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="logo" prop="logo">
          <el-input v-model="project.thumbnail" placeholder="请输入logo地址" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="project.srcUrl" placeholder="请输入地址" />
        </el-form-item>
        <el-form-item label="展示地址" prop="displayUrl">
          <el-input v-model="project.displayUrl" placeholder="请输入展示地址" />
        </el-form-item>
        <el-form-item label="开始日期" prop="beginDate">
          <el-date-picker v-model="project.beginDate" type="date" placeholder="选择日期" />
        </el-form-item>
        <el-form-item label="公开" prop="status">
          <el-switch v-model="project.status" active-value="0" inactive-value="1" />
        </el-form-item>
        <el-form-item label="置顶" prop="isTop">
          <el-switch v-model="project.isTop" active-value="1" inactive-value="0" />
        </el-form-item>
        <el-form-item label="排序" prop="orderNum">
          <el-input v-model="project.orderNum" placeholder="请输入排序" />
        </el-form-item>
        <el-form-item label="项目介绍">
          <MdEditor v-model="project.content" />
        </el-form-item>
      </el-form>
      <template v-slot:footer>
        <div class="dialog-footer">
          <el-button @click="handleCancel">取 消</el-button>
          <el-button type="primary" @click="handleSubmit">确 定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Search, Plus, Delete, Edit } from '@element-plus/icons-vue'
import type { FormRules } from 'element-plus'
import { toggleStatus } from '@/utils/toggleStatus'
import type { ProjectAddParams, ProjectListData, ProjectPageParams } from '@/api/project/type'
import {
  reqProjectAdd,
  reqProjectDelete,
  reqProjectGetById,
  reqProjectPage,
  reqProjectUpdate,
} from '@/api/project'

const queryParams = ref<ProjectPageParams>({
  pageNum: 1,
  pageSize: 10,
})
const projectList = ref<ProjectListData[]>([])
const total = ref<number>(0)
const loading = ref<boolean>(false)
const open = ref<boolean>(false)
const title = ref<string>('')
const project = ref<ProjectAddParams>({})
const rules = ref<FormRules>({})
const statusLoading = ref<boolean>(false)

const getProjectPage = async (page: number = 1) => {
  loading.value = true
  queryParams.value.pageNum = page
  const res = await reqProjectPage(queryParams.value)
  projectList.value = res.data.rows
  total.value = res.data.total
  loading.value = false
}

const getProjectDetail = async (id: number) => {
  const res = await reqProjectGetById(id)
  return res.data
}

const handleAdd = () => {
  project.value.status = '0'
  project.value = {}
  open.value = true
  title.value = '新增'
}

const handleUpdate = async (row: ProjectListData) => {
  project.value = await getProjectDetail(row.id)
  open.value = true
  title.value = '修改' + row.title
}

const handleDelete = (row: ProjectListData) => {
  ElMessageBox.confirm('删除?', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async () => {
    await reqProjectDelete(row.id)
    await getProjectPage(queryParams.value.pageNum)
  })
}

const handleChangeStatus = async (row: ProjectListData) => {
  statusLoading.value = true
  try {
    await reqProjectUpdate(row.id, { status: row.status, isTop: row.isTop })
  } catch {
    toggleStatus(row)
  } finally {
    statusLoading.value = false
  }
}

const handleSubmit = async () => {
  if (project.value.id) {
    await reqProjectUpdate(project.value.id, project.value)
  } else {
    await reqProjectAdd(project.value)
    ElMessage.success('新增成功')
  }
  open.value = false
  getProjectPage()
}

const handleCancel = () => {
  project.value = {}
  open.value = false
}

onMounted(() => {
  getProjectPage()
})
</script>

<style scoped lang="scss"></style>
