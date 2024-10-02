<template>
  <div class="app-container">
    <el-form ref="queryForm" :model="queryParams" class="query-form" size="small">
      <el-input
        v-model="queryParams.name"
        placeholder="分类名"
        clearable
        @keyup.enter="getCategoryList()"
      />
      <el-select v-model="queryParams.status" placeholder="状态" clearable>
        <el-option key="0" label="正常" value="0" />
        <el-option key="1" label="禁用" value="1" />
      </el-select>
      <el-button type="primary" :icon="Search" @click="getCategoryList()">搜索</el-button>
    </el-form>

    <el-row :gutter="10">
      <el-col :span="1.5">
        <el-button type="primary" plain :icon="Plus" size="mini" @click="handleAdd">新增</el-button>
      </el-col>

      <el-col :span="1.5">
        <el-button
          type="warning"
          plain
          :icon="Download"
          :loading="exportLoading"
          @click="handleExport"
        >
          导出
        </el-button>
      </el-col>
    </el-row>

    <el-table v-loading="loading" :data="categoryList">
      <el-table-column label="id" align="center" prop="id" />
      <el-table-column label="分类名" align="center" prop="name" />
      <el-table-column label="描述" align="center" prop="description" />
      <el-table-column
        label="操作"
        align="center"
        class-name="small-padding fixed-width"
        fixed="right"
      >
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
      layout="sizes, prev, pager, next, jumper, ->, total"
      :total="total"
      :page-sizes="[5, 10, 20, 30, 40]"
      v-model:current-page="queryParams.pageNum"
      @current-change="getCategoryList"
      @size-change="getCategoryList"
    />

    <!-- 添加或修改分类对话框 -->
    <el-dialog :title="title" v-model="open" class="dialog-form">
      <el-form ref="form" :model="category" :rules="rules" label-width="auto" label-position="top">
        <el-form-item label="分类名" prop="name">
          <el-input v-model="category.name" placeholder="请输入分类名" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="category.description" type="textarea" placeholder="请输入内容" />
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
import { Search, Plus, Delete, Download, Edit } from '@element-plus/icons-vue'
import {
  reqCategoryAdd,
  reqCategoryDelete,
  reqCategoryPage,
  reqCategoryUpdate,
  reqExportToXlsx,
} from '@/api/content/category'
import type {
  CategoryAddParams,
  CategoryListData,
  CategoryPageParams,
} from '@/api/content/category/type'

const categoryList = ref<CategoryListData[]>()
const queryParams = ref<CategoryPageParams>({
  pageNum: 1,
  pageSize: 10,
})
const total = ref<number>(0)
const open = ref<boolean>(false)
const category = ref<CategoryAddParams>({})
const title = ref<string>('')
const rules = ref({})
const exportLoading = ref<boolean>(false)
const loading = ref<boolean>(false)

const getCategoryList = async (page: number = 1) => {
  loading.value = true
  queryParams.value.pageNum = page
  const res = await reqCategoryPage(queryParams.value)
  categoryList.value = res.data.rows
  total.value = res.data.total
  loading.value = false
}

const handleAdd = () => {
  category.value = {}
  open.value = true
  title.value = '新增'
}

const handleExport = async () => {
  exportLoading.value = true
  await reqExportToXlsx()
  exportLoading.value = false
}

const handleUpdate = (row: CategoryListData) => {
  category.value = row
  open.value = true
  title.value = '修改' + row.name
}

const handleSubmit = async () => {
  if (category.value.id) {
    await reqCategoryUpdate(category.value.id, category.value)
  } else {
    await reqCategoryAdd(category.value)
    ElMessage.success('新增成功')
  }
  open.value = false
  await getCategoryList(queryParams.value.pageNum)
}

const handleCancel = () => {
  category.value = {}
  open.value = false
}

const handleDelete = async (row: CategoryListData) => {
  ElMessageBox.confirm('删除?', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async () => {
    await reqCategoryDelete(row.id)
    await getCategoryList(queryParams.value.pageNum)
  })
}

onMounted(async () => {
  await getCategoryList()
})
</script>

<style scoped lang="scss"></style>
