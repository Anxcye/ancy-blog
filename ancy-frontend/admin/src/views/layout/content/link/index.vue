<template>
  <div class="app-container">
    <el-form ref="queryForm" :model="queryParams" :inline="true" size="small">
      <el-row :gutter="10" style="margin-bottom: 10px">
        <el-col :span="10">
          <el-input
            v-model="queryParams.name"
            placeholder="名称"
            clearable
            @keyup.enter="getLinkPage()"
          />
        </el-col>

        <el-col :span="6">
          <el-select v-model="queryParams.status" placeholder="公开?" clearable>
            <el-option :key="'0'" label="公开" :value="'0'" />
            <el-option :key="'1'" label="隐藏" :value="'1'" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-button type="primary" :icon="Search" @click="getLinkPage()">搜索</el-button>
        </el-col>
      </el-row>
    </el-form>

    <el-row :gutter="10">
      <el-col :span="1.5">
        <el-button type="primary" plain :icon="Plus" @click="handleAdd">新增</el-button>
      </el-col>
    </el-row>

    <el-table v-loading="loading" :data="linkList">
      <el-table-column label="ID" align="center" width="50" prop="id" />
      <el-table-column label="名称" align="center" prop="name" />
      <el-table-column label="描述" align="center" prop="description" />
      <el-table-column label="logo" align="center" width="100" prop="logo" type="img">
        <template v-slot="scope">
          <el-image style="width: 100px; height: 100px" :src="scope.row.logo" fit="fill" />
        </template>
      </el-table-column>
      <el-table-column prop="address" label="地址" align="center">
        <template v-slot="scope">
          <el-link :href="scope.row.address" target="_blank" style="word-break: break-all">
            {{ scope.row.address }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="公开" align="center">
        <template v-slot="scope">
          <el-switch
            v-model="scope.row.status"
            active-value="0"
            inactive-value="1"
            @change="handleChangeStatus(scope.row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" fixed="right">
        <template v-slot="scope">
          <el-button type="text" :icon="Edit" @click="handleUpdate(scope.row)">修改</el-button>
          <el-button type="text" :icon="Delete" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:page-size="queryParams.pageSize"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total"
      :page-sizes="[5, 10, 20, 30, 40]"
      v-model:current-page="queryParams.pageNum"
      @current-change="getLinkPage"
      @size-change="getLinkPage()"
    />

    <el-dialog :title="title" v-model="open" width="500px" append-to-body>
      <el-form ref="form" :rules="rules" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="link.name" placeholder="请输入名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="link.description" type="textarea" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="logo" prop="logo">
          <el-input v-model="link.logo" placeholder="请输入logo地址" />
        </el-form-item>
        <el-form-item label="地址" prop="address">
          <el-input v-model="link.address" placeholder="请输入地址" />
        </el-form-item>
        <el-form-item label="公开" prop="status">
          <el-switch v-model="link.status" active-value="0" inactive-value="1" />
        </el-form-item>
      </el-form>
      <template v-slot:footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="handleSubmit">确 定</el-button>
          <el-button @click="handleCancel">取 消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reqLinkAdd, reqLinkDelete, reqLinkPage, reqLinkUpdate } from '@/api/content/link'
import type { LinkAddParams, LinkListData, LinkPageParams } from '@/api/content/link/type'
import { onMounted, ref } from 'vue'
import { Search, Plus, Delete, Edit } from '@element-plus/icons-vue'
import type { FormRules } from 'element-plus'

const queryParams = ref<LinkPageParams>({
  pageNum: 1,
  pageSize: 10,
})
const linkList = ref<LinkListData[]>([])
const total = ref<number>(0)
const loading = ref<boolean>(false)
const open = ref<boolean>(false)
const title = ref<string>('')
const link = ref<LinkAddParams>({})
const rules = ref<FormRules>({})
const statusLoading = ref<boolean>(false)

const getLinkPage = async (page: number = 1) => {
  loading.value = true
  queryParams.value.pageNum = page
  const res = await reqLinkPage(queryParams.value)
  linkList.value = res.data.rows
  total.value = res.data.total
  loading.value = false
}

const handleAdd = () => {
  link.value = {}
  open.value = true
  title.value = '新增'
}

const handleUpdate = (row: LinkListData) => {
  link.value = row
  open.value = true
  title.value = '修改' + row.name
}

const handleDelete = (row: LinkListData) => {
  ElMessageBox.confirm('删除?', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async () => {
    await reqLinkDelete(row.id)
    await getLinkPage()
  })
}

const handleChangeStatus = async (row: LinkListData) => {
  statusLoading.value = true
  await reqLinkUpdate(row.id, { status: row.status })
  statusLoading.value = false
}

const handleSubmit = async () => {
  if (link.value.id) {
    await reqLinkUpdate(link.value.id, link.value)
  } else {
    await reqLinkAdd(link.value)
    ElMessage.success('新增成功')
  }
  open.value = false
  getLinkPage()
}

const handleCancel = () => {
  open.value = false
}

onMounted(() => {
  getLinkPage()
})
</script>

<style scoped lang="scss"></style>
