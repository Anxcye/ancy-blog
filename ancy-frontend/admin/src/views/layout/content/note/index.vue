<template>
  <div class="app-container">
    <el-form ref="queryForm" :model="queryParams" class="query-form" size="small">
      <el-form-item prop="content">
        <el-input
          v-model="queryParams.content"
          placeholder="内容"
          clearable
          @keyup.enter="getNotePage()"
        />
      </el-form-item>
      <el-form-item prop="status">
        <el-select v-model="queryParams.status" placeholder="公开?" clearable>
          <el-option key="0" label="公开" value="0" />
          <el-option key="1" label="隐藏" value="1" />
        </el-select>
      </el-form-item>

      <el-form-item>
        <el-button type="primary" :icon="Search" @click="getNotePage()">搜索</el-button>
      </el-form-item>
    </el-form>

    <el-row>
      <el-button type="primary" plain :icon="Plus" @click="handleAdd">新增</el-button>
    </el-row>

    <el-table :data="NoteList">
      <el-table-column prop="id" width="80" label="ID" align="center" />
      <el-table-column
        prop="content"
        width="300"
        show-overflow-tooltip
        label="内容"
        align="center"
      />
      <el-table-column prop="viewCount" label="浏览量" align="center" />
      <el-table-column prop="status" label="公开" align="center">
        <template v-slot="scope">
          <el-switch
            v-model="scope.row.status"
            active-value="0"
            inactive-value="1"
            :loading="statusLoading"
            @change="handleStatusChange(scope.row)"
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
            @change="handleStatusChange(scope.row)"
          />
        </template>
      </el-table-column>
      <el-table-column prop="isComment" label="允许评论" align="center">
        <template v-slot="scope">
          <el-switch
            v-model="scope.row.isComment"
            active-value="1"
            inactive-value="0"
            :loading="statusLoading"
            @change="handleStatusChange(scope.row)"
          />
        </template>
      </el-table-column>
      <el-table-column prop="orderNum" label="排序" align="center" />

      <el-table-column label="时间" align="center" width="250" prop="createTime">
        <template v-slot="scope">
          <div>创建时间: {{ scope.row.createTime }}</div>
          <div v-if="scope.row.updateTime">更新时间: {{ scope.row.updateTime }}</div>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="170px" fixed="right">
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
      @current-change="getNotePage"
      @size-change="getNotePage()"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Plus, Delete, Edit, Search } from '@element-plus/icons-vue'
import router from '@/router'
import type { NotePageData, NotePageParams } from '@/api/note/type'
import { reqNoteDelete, reqNotePage, reqNoteUpdate } from '@/api/note'

const queryParams = ref<NotePageParams>({
  pageNum: 1,
  pageSize: 10,
})
const NoteList = ref<NotePageData[]>([])
const total = ref(0)
const statusLoading = ref(false)

const getNotePage = async (page: number = 1) => {
  queryParams.value.pageNum = page
  const res = await reqNotePage(queryParams.value)
  NoteList.value = res.data.rows
  total.value = res.data.total
}

onMounted(async () => {
  await getNotePage()
})

const handleAdd = () => {
  router.push('/write-note')
}

const handleDelete = (row: NotePageData) => {
  ElMessageBox.confirm('删除?', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async () => {
    await reqNoteDelete(row.id)
    await getNotePage(queryParams.value.pageNum)
  })
}

const handleUpdate = (a: NotePageData) => {
  router.push({
    path: '/write-note',
    query: {
      id: a.id,
    },
  })
}

const handleStatusChange = async (row: NotePageData) => {
  statusLoading.value = true
  try {
    await reqNoteUpdate(row.id, {
      status: row.status,
      isTop: row.isTop,
      isComment: row.isComment,
    })
  } catch {
    getNotePage(queryParams.value.pageNum)
  } finally {
    statusLoading.value = false
  }
}
</script>

<style scoped lang="scss"></style>
