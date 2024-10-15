<template>
  <div class="app-container">
    <el-form ref="queryForm" :model="queryParams" class="query-form" size="small">
      <el-input
        v-model="queryParams.name"
        placeholder="标签名"
        clearable
        @keyup.enter="getTagList()"
      />
      <el-button type="primary" :icon="Search" @click="getTagList()">搜索</el-button>
    </el-form>

    <el-row :gutter="10">
      <el-col :span="1.5">
        <el-button type="primary" plain :icon="Plus" @click="handleAdd">新增</el-button>
      </el-col>
    </el-row>

    <el-table v-loading="loading" :data="tagList">
      <el-table-column label="ID" align="center" width="50" prop="id" />
      <el-table-column label="标签名" align="center" prop="name" />
      <el-table-column label="备注" align="center" prop="remark" />
      <el-table-column label="操作" align="center" width="170" fixed="right">
        <template v-slot="scope">
          <el-button size="small" type="text" :icon="Edit" @click="handleUpdate(scope.row)">
            修改
          </el-button>
          <el-button size="small" type="text" :icon="Delete" @click="handleDelete(scope.row)">
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
      @current-change="getTagList"
      @size-change="getTagList"
    />

    <!-- 添加或修改分类对话框 -->
    <el-dialog :title="title" v-model="open" class="dialog-form">
      <el-form ref="form" :model="tag" :rules="rules" label-width="auto" label-position="top">
        <el-form-item prop="name" label="标签名">
          <el-input v-model="tag.name" placeholder="标签名" />
        </el-form-item>
        <el-form-item prop="remark" label="备注">
          <el-input v-model="tag.remark" type="textarea" placeholder="备注" />
        </el-form-item>
      </el-form>
      <template v-slot:footer>
        <div class="dialog-footer">
          <el-button @click="hancleCancel">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reqTagAdd, reqTagDelete, reqTagPage, reqTagUpdate } from '@/api/content/tag'
import type { TagAddData, TagListData, TagPageParams } from '@/api/content/tag/type'
import { onMounted, ref } from 'vue'
import { Search, Delete, Plus, Edit } from '@element-plus/icons-vue'

const queryParams = ref<TagPageParams>({
  pageNum: 1,
  pageSize: 10,
})
const loading = ref(false)
const tagList = ref<TagListData[]>([])
const total = ref<number>(0)
const title = ref<string>('')
const open = ref<boolean>(false)
const tag = ref<TagAddData>({})
const rules = ref<any>({})

const getTagList = async (page: number = 1) => {
  loading.value = true
  queryParams.value.pageNum = page
  const res = await reqTagPage(queryParams.value)
  tagList.value = res.data.rows
  total.value = res.data.total
  loading.value = false
}

const handleAdd = () => {
  title.value = '新增'
  open.value = true
}

const handleUpdate = (row: TagListData) => {
  tag.value = row
  title.value = '修改'
  open.value = true
}

const handleDelete = (row: TagListData) => {
  ElMessageBox.confirm('删除?', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async () => {
    await reqTagDelete(row.id)
    await getTagList()
  })
}

const handleSubmit = async () => {
  if (tag.value.id) {
    await reqTagUpdate(tag.value.id, tag.value)
  } else {
    await reqTagAdd(tag.value)
    ElMessage.success('新增成功')
  }
  open.value = false
  getTagList()
}

const hancleCancel = () => {
  open.value = false
  tag.value = {}
}

onMounted(() => {
  getTagList()
})
</script>

<style scoped lang="scss"></style>
