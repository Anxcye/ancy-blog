<template>
  <div class="app-container">
    <el-form ref="queryForm" :model="queryParams" class="query-form" size="small">
      <el-form-item prop="title">
        <el-input
          v-model="queryParams.title"
          placeholder="标题"
          clearable
          @keyup.enter="getArticlePage()"
        />
      </el-form-item>
      <el-form-item prop="summary">
        <el-input
          v-model="queryParams.summary"
          placeholder="摘要"
          clearable
          @keyup.enter="getArticlePage()"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :icon="Search" @click="getArticlePage()">搜索</el-button>
      </el-form-item>
    </el-form>

    <el-row>
      <el-button type="primary" plain :icon="Plus" @click="handleAdd">新增</el-button>
    </el-row>

    <el-table :data="articleList">
      <el-table-column prop="id" width="80" label="ID" align="center" />
      <el-table-column prop="title" width="200" label="标题" align="center" />
      <el-table-column prop="summary" width="200" label="摘要" align="center" />
      <el-table-column prop="categoryName" label="分类" align="center" />
      <el-table-column prop="tags" label="标签" align="center">
        <template v-slot="scope">
          <el-tag v-for="tag in scope.row.tags" :key="tag.id" size="small">
            {{ tag.name }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="viewCount" label="浏览量" align="center" />
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

      <el-table-column prop="type" label="挂在首页下" align="center">
        <template v-slot="scope">
          <el-switch
            v-model="scope.row.type"
            :active-value="1"
            :inactive-value="0"
            :loading="statusLoading"
            @change="handleStatusChange(scope.row)"
          />
        </template>
      </el-table-column>

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
      @current-change="getArticlePage"
      @size-change="getArticlePage()"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Plus, Delete, Edit, Search } from '@element-plus/icons-vue'
import { reqArticleDelete, reqArticlePage, reqArticleUpdate } from '@/api/content/article'
import type { ArticlePageData, ArticlePageParams } from '@/api/content/article/type'
import router from '@/router'

const queryParams = ref<ArticlePageParams>({
  pageNum: 1,
  pageSize: 10,
  title: '',
  summary: '',
})
const articleList = ref<ArticlePageData[]>([])
const total = ref(0)
const statusLoading = ref(false)

const getArticlePage = async (page: number = 1) => {
  queryParams.value.pageNum = page
  const res = await reqArticlePage(queryParams.value)
  articleList.value = res.data.rows
  total.value = res.data.total
}

onMounted(async () => {
  await getArticlePage()
})

const handleAdd = () => {
  router.push('/write')
}

const handleDelete = (article: ArticlePageData) => {
  ElMessageBox.prompt('请键入标题"' + article.title + '"以确认删除', {
    confirmButtonText: '删除',
    cancelButtonText: '取消',
    inputPlaceholder: article.title,
  })
    .then(async function (value) {
      if (value.value === article.title) {
        await reqArticleDelete(article.id)
        getArticlePage()
      } else {
        ElMessage.error('输入不一致，删除失败')
      }
    })
    .then(() => {
      getArticlePage()
    })
    .catch(() => {
      ElMessage.error('取消删除')
    })
}

const handleUpdate = (a: ArticlePageData) => {
  router.push({
    path: '/write',
    query: {
      id: a.id,
    },
  })
}

const handleStatusChange = async (article: ArticlePageData) => {
  statusLoading.value = true
  try {
    await reqArticleUpdate(article.id, {
      status: article.status,
      isTop: article.isTop,
      type: article.type,
    })
  } catch {
    getArticlePage(queryParams.value.pageNum)
  } finally {
    statusLoading.value = false
  }
}
</script>

<style scoped lang="scss"></style>
