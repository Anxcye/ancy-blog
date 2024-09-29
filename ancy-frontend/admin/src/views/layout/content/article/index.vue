<template>
  <div class="container">
    <el-form
      ref="queryForm"
      :model="queryParams"
      :inline="true"
      label-width="68px"
      class="search-form"
    >
      <el-form-item prop="title">
        <el-input
          v-model="queryParams.title"
          placeholder="标题"
          clearable
          size="small"
          style="max-width: 200px"
          @keyup.enter="getArticlePage()"
        />
      </el-form-item>
      <el-form-item prop="summary">
        <el-input
          v-model="queryParams.summary"
          placeholder="摘要"
          clearable
          size="small"
          style="max-width: 200px"
          @keyup.enter="getArticlePage()"
        />
      </el-form-item>
      <el-form-item>
        <el-button
          type="primary"
          :icon="Search"
          size="small"
          @click="getArticlePage()"
        >
          搜索
        </el-button>
      </el-form-item>
    </el-form>

    <el-row :gutter="10">
      <el-button type="primary" plain :icon="Plus" @click="handleAdd">
        新增
      </el-button>
    </el-row>

    <el-table :data="articleList">
      <el-table-column prop="id" width="80" label="ID" align="center" />
      <el-table-column prop="title" width="200" label="标题" align="center" />
      <el-table-column prop="summary" width="200" label="摘要" align="center" />
      <el-table-column prop="createTime" label="创建时间" align="center" />

      <el-table-column
        label="操作"
        align="center"
        width="170px"
        class-name="small-padding fixed-width"
        fixed="right"
      >
        <template v-slot="scope">
          <el-button
            size="small"
            type="text"
            :icon="Edit"
            @click="handleUpdate(scope.row)"
          >
            修改
          </el-button>
          <el-button
            size="small"
            type="text"
            :icon="Delete"
            @click="handleDelete(scope.row)"
          >
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
import { reqArticleDelete, reqArticlePage } from '@/api/content/article'
import type {
  ArticlePageData,
  ArticlePageParams,
} from '@/api/content/article/type'

const queryParams = ref<ArticlePageParams>({
  pageNum: 1,
  pageSize: 3,
  title: '',
  summary: '',
})

const articleList = ref<ArticlePageData[]>([])
const total = ref(0)

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
  console.log('handleAdd')
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

const handleUpdate = (a) => {
  console.log('handleUpdate', a)
}
</script>

<style scoped lang="scss">
.container {
  width: 100%;
}
</style>
