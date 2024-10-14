<template>
  <div class="app-container">
    <el-form ref="queryForm" :model="queryParams" class="query-form" size="small">
      <el-input
        v-model="queryParams.articleId"
        placeholder="文章id"
        clearable
        @keyup.enter="getCommentList()"
      />
      <el-input
        v-model="queryParams.nickname"
        placeholder="昵称"
        clearable
        @keyup.enter="getCommentList()"
      />
      <el-input
        v-model="queryParams.email"
        placeholder="E-mail"
        clearable
        @keyup.enter="getCommentList()"
      />
      <el-input
        v-model="queryParams.content"
        placeholder="内容"
        clearable
        @keyup.enter="getCommentList()"
      />
      <el-select v-model="queryParams.status" placeholder="状态" clearable>
        <el-option key="0" label="公开" value="0" />
        <el-option key="1" label="隐藏" value="1" />
      </el-select>
      <el-button type="primary" :icon="Search" @click="getCommentList()">搜索</el-button>
    </el-form>

    <el-row :gutter="10">
      <el-col :span="1.5">
        <el-button type="primary" plain :icon="Plus" size="mini" @click="handleAdd">新增</el-button>
      </el-col>
    </el-row>

    <el-table v-loading="loading" :data="commentList">
      <el-table-column label="id" align="center" prop="id" />
      <el-table-column label="文章" align="center" prop="articleId">
        <template #default="scope">
          <el-link
            :href="`${baseInfoStore.baseInfo.address}/${getCommentArticleUri(scope.row)}`"
            target="_blank"
          >
            {{ getCommentArticleUri(scope.row) }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column label="昵称" align="center" prop="nickname" />
      <el-table-column
        label="内容"
        align="center"
        width="200"
        show-overflow-tooltip
        prop="content"
      />
      <el-table-column label="类型" align="center">
        <template #default="scope">
          {{ scope.row.parentId === -1 ? '评论' : `回复${scope.row.toCommentNickname ?? ''}` }}
        </template>
      </el-table-column>
      <el-table-column label="状态" align="center" prop="status">
        <template #default="scope">
          <el-switch
            v-model="scope.row.status"
            :loading="statusLoading"
            active-value="0"
            inactive-value="1"
            @change="handleStatusChange(scope.row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="ua" align="center" width="200" show-overflow-tooltip prop="ua" />
      <el-table-column label="ip" align="center" width="200" show-overflow-tooltip prop="ip" />
      <el-table-column label="E-mail" align="center" prop="email" />
      <el-table-column label="头像" align="center" prop="avatar">
        <template #default="scope">
          <el-image preview-teleported :src="scope.row.avatar" style="width: 50px; height: 50px" />
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="100" fixed="right">
        <template v-slot="scope">
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
      @current-change="getCommentList"
      @size-change="getCommentList"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Search, Plus, Delete } from '@element-plus/icons-vue'
import type { CommentListData, CommentPageParams } from '@/api/comment/type'
import { reqCommentDelete, reqCommentPage, reqCommentUpdate } from '@/api/comment'
import { useBaseInfoStore } from '@/stores/modules/baseInfo'

const baseInfoStore = useBaseInfoStore()
const commentList = ref<CommentListData[]>()
const queryParams = ref<CommentPageParams>({
  pageNum: 1,
  pageSize: 10,
})
const total = ref<number>(0)
const statusLoading = ref<boolean>(false)
const loading = ref<boolean>(false)

const getCommentArticleUri = (row: CommentListData) => {
  return `/${row.type === '0' ? 'article' : 'note'}/${row.articleId}`
}

const getCommentList = async (page: number = 1) => {
  loading.value = true
  queryParams.value.pageNum = page
  const res = await reqCommentPage(queryParams.value)
  commentList.value = res.data.rows
  total.value = res.data.total
  loading.value = false
}

const handleAdd = () => {
  ElMessageBox.confirm('登录前台以新增评论', '', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async () => {
    window.open(`${baseInfoStore.baseInfo.address}`, '_blank')
  })
}

const handleStatusChange = async (row: CommentListData) => {
  if (statusLoading.value) return
  statusLoading.value = true
  await reqCommentUpdate(row.id, { status: row.status }).finally(() => {
    statusLoading.value = false
  })
}

const handleDelete = async (row: CommentListData) => {
  ElMessageBox.confirm('删除?', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async () => {
    await reqCommentDelete(row.id)
    await getCommentList(queryParams.value.pageNum)
  })
}

onMounted(async () => {
  await getCommentList()
})
</script>

<style scoped lang="scss"></style>
