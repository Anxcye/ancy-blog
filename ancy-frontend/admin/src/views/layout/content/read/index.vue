<template>
  <div class="app-container">
    <el-form ref="queryForm" :model="queryParams" class="query-form" size="small">
      <el-input
        v-model="queryParams.author"
        placeholder="作者"
        clearable
        @keyup.enter="getReadPage()"
      />
      <el-input
        v-model="queryParams.source"
        placeholder="来源"
        clearable
        @keyup.enter="getReadPage()"
      />
      <el-input
        v-model="queryParams.content"
        placeholder="内容"
        clearable
        @keyup.enter="getReadPage()"
      />
      <el-select v-model="queryParams.addFrom" placeholder="添加来源" clearable>
        <el-option key="0" label="手动" value="0" />
        <el-option key="1" label="安读" value="1" />
      </el-select>
      <el-button type="primary" :icon="Search" @click="getReadPage()">搜索</el-button>
    </el-form>

    <el-row :gutter="10">
      <el-col :span="1.5">
        <el-button type="primary" plain :icon="Plus" @click="handleAdd">新增</el-button>
      </el-col>
    </el-row>

    <el-table v-loading="loading" :data="readList">
      <el-table-column label="ID" align="center" width="50" prop="id" />
      <el-table-column label="作者" align="center" width="150" prop="author" />
      <el-table-column label="来源" align="center" width="150" prop="source" />

      <el-table-column label="内容" align="center" show-overflow-tooltip prop="content" />

      <el-table-column label="添加来源" align="center" width="100" prop="addFrom">
        <template v-slot="scope">
          <el-tag v-if="scope.row.addFrom === 0">手动</el-tag>
          <el-tag v-else type="success">安读</el-tag>
        </template>
      </el-table-column>

      <el-table-column label="时间" align="center" width="250" prop="createTime">
        <template v-slot="scope">
          <div>创建时间: {{ scope.row.createTime }}</div>
          <div v-if="scope.row.updateTime">更新时间: {{ scope.row.updateTime }}</div>
        </template>
      </el-table-column>

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
      @current-change="getReadPage"
      @size-change="getReadPage()"
    />

    <el-dialog :title="title" v-model="open" class="dialog-form">
      <el-form ref="form" :rules="rules" label-width="auto" label-position="top">
        <el-form-item label="作者" prop="author">
          <el-input v-model="read.author" placeholder="作者" />
        </el-form-item>
        <el-form-item label="来源" prop="source">
          <el-input v-model="read.source" placeholder="来源 可能是一本书名" />
        </el-form-item>
        <el-form-item label="内容" prop="content">
          <el-input
            v-model="read.content"
            placeholder="说了什么呢？"
            type="textarea"
            auto-size
            :autosize="{ minRows: 2, maxRows: 7 }"
          />
        </el-form-item>
        <el-form-item label="添加来源" prop="addFrom">
          <el-select v-model="read.addFrom" placeholder="手动" disabled default-first-option>
            <el-option :key="0" label="手动" :value="0" />
            <el-option :key="1" label="安读" :value="1" />
          </el-select>
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
import type { ReadAddParams, ReadListData, ReadPageParams } from '@/api/read/type'
import { reqReadAdd, reqReadDelete, reqReadPage, reqReadUpdate } from '@/api/read'

const queryParams = ref<ReadPageParams>({
  pageNum: 1,
  pageSize: 10,
})
const readList = ref<ReadListData[]>([])
const total = ref<number>(0)
const loading = ref<boolean>(false)
const open = ref<boolean>(false)
const title = ref<string>('')
const read = ref<ReadAddParams>({})
const rules = ref<FormRules>({})

const getReadPage = async (page: number = 1) => {
  loading.value = true
  queryParams.value.pageNum = page
  const res = await reqReadPage(queryParams.value)
  readList.value = res.data.rows
  total.value = res.data.total
  loading.value = false
}

const handleAdd = () => {
  read.value = {}
  open.value = true
  title.value = '新增'
}

const handleUpdate = (row: ReadListData) => {
  read.value = { ...row }
  open.value = true
  title.value = '修改' + row.author
}

const handleDelete = (row: ReadListData) => {
  ElMessageBox.confirm('删除?', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async () => {
    await reqReadDelete(row.id)
    await getReadPage(queryParams.value.pageNum)
  })
}

const handleSubmit = async () => {
  if (read.value.id) {
    await reqReadUpdate(read.value.id, read.value)
  } else {
    read.value.addFrom = 0
    await reqReadAdd(read.value)
    ElMessage.success('新增成功')
  }
  open.value = false
  getReadPage()
}

const handleCancel = () => {
  read.value = {}
  open.value = false
}

onMounted(() => {
  getReadPage()
})
</script>

<style scoped lang="scss"></style>
