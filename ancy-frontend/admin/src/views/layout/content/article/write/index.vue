<template>
  <div class="app-container">
    <el-form ref="form" :model="article" label-width="90px">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-input v-model="article.title" placeholder="文章标题" maxlength="30" />
        </el-col>
        <el-col :span="6">
          <el-select
            v-model="article.categoryId"
            placeholder="选择分类，支持创建"
            allow-create
            filterable
            default-first-option
          >
            <el-option
              v-for="category in categoryList"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select
            v-model="article.tags"
            placeholder="选择标签，支持创建"
            multiple
            allow-create
            filterable
            default-first-option
          >
            <el-option v-for="tag in tagList" :key="tag.id" :label="tag.name" :value="tag.id" />
          </el-select>
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-top: 10px">
        <el-col :span="12">
          <el-input
            v-model="article.summary"
            type="textarea"
            :autosize="{ minRows: 3, maxRows: 5 }"
            placeholder="文章摘要"
          />
        </el-col>
        <el-col :span="12">
          <el-upload
            v-if="!article.thumbnail"
            list-type="picture"
            drag
            name="img"
            action="upload"
            :limit="1"
            :http-request="handleUpload"
            :show-file-list="false"
          >
            <div class="el-upload__text">
              推拽以添加缩略图，或
              <em>点击上传</em>
              <div class="el-upload__tip">jpg/png 不超过500kb</div>
            </div>
          </el-upload>
          <div v-else class="img-box">
            <img :src="article.thumbnail" alt="" />
            <el-button type="text" @click="article.thumbnail = undefined">删除</el-button>
          </div>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin: 10px 0">
        <el-col :span="7">
          <el-form-item label="评论" label-width="40px">
            <el-radio-group v-model="article.isComment">
              <el-radio key="1" label="1">允许</el-radio>
              <el-radio key="0" label="0">禁止</el-radio>
            </el-radio-group>
            {{ article.isComment }}
          </el-form-item>
        </el-col>
        <el-col :span="7" style="text-align: right">
          <el-form-item label="置顶" label-width="40px">
            <el-radio-group v-model="article.isTop">
              <el-radio key="1" label="1">是</el-radio>
              <el-radio key="0" label="0">否</el-radio>
              {{ article.isTop }}
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="10" style="text-align: right">
          <el-button v-if="!aId" type="info" @click="handleDraft">保存到草稿箱</el-button>
          <el-button type="primary" @click="handleSubmit">
            {{ aId ? '更新' : '发布' }}
          </el-button>
        </el-col>
      </el-row>
      <el-row>
        <MdEditor v-model="article.content" />
      </el-row>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import MdEditor from './MdEditor.vue'
import type { CategoryListData } from '@/api/content/category/type'
import type { ArticleAddParams } from '@/api/content/article/type'
import { reqCategoryAdd, reqCategoryList } from '@/api/content/category'
import type { TagListData } from '@/api/content/tag/type'
import { reqTagAdd, reqTagList } from '@/api/content/tag'
import { reqUpload } from '@/api/common'
import type { UploadRequestOptions } from 'element-plus'
import { useRoute, type LocationQueryValue } from 'vue-router'
import { reqArticleAdd, reqArticleGetById, reqArticleUpdate } from '@/api/content/article'

const route = useRoute()
const article = ref<ArticleAddParams>({})
const categoryList = ref<CategoryListData[]>([])
const tagList = ref<TagListData[]>([])
const aId = ref<number | null>(null)

watch(
  () => route.query.id,
  async (newVal: LocationQueryValue | LocationQueryValue[]) => {
    aId.value = Number(newVal)
    if (aId.value) {
      await getArticleById(aId.value)
    } else {
      article.value = {}
    }
  },
  {
    immediate: true,
  },
)

watch(
  article,
  (newVal) => {
    localStorage.setItem('write-content', JSON.stringify(newVal))
  },
  {
    deep: true,
  },
)

watch(
  () => article.value.categoryId,
  async (newVal) => {
    if (!newVal) return
    if (!categoryList.value.find((item) => item.id === newVal)) {
      const res = await reqCategoryAdd({
        name: newVal.toString(),
      })
      await getCategoryList()
      article.value.categoryId = res.data
    }
  },
)

watch(
  () => article.value.tags,
  async (newVal) => {
    if (!newVal) return
    for (const tag of newVal!) {
      if (!tagList.value.find((item) => item.id === tag)) {
        const res = await reqTagAdd({
          name: tag!.toString(),
        })
        await getTagList()
        article.value.tags!.splice(article.value.tags!.indexOf(tag), 1)
        article.value.tags!.push(res.data)
      }
    }
  },
)

const getCategoryList = async () => {
  const res = await reqCategoryList()
  categoryList.value = res.data
}

const getTagList = async () => {
  const res = await reqTagList()
  tagList.value = res.data
}

const getArticleById = async (id: number) => {
  const res = await reqArticleGetById(id)
  article.value = {
    ...res.data,
    tags: res.data.tags.map((item) => item.id),
  }
}

const handleSubmit = async () => {
  if (aId.value) {
    await reqArticleUpdate(aId.value, article.value)
  } else {
    await reqArticleAdd(article.value)
  }
}

const handleDraft = async () => {
  article.value.status = '1'
  await reqArticleAdd(article.value)
}

const handleUpload = async (img: UploadRequestOptions) => {
  const res = await reqUpload(img.file)
  article.value.thumbnail = res.data
  return res
}

onMounted(async () => {
  await getCategoryList()
  await getTagList()
  if (aId.value) {
    await getArticleById(aId.value)
  }

  if (!aId.value && localStorage.getItem('write-content')) {
    ElMessageBox.confirm('是否恢复上次编辑内容', '提示', {
      confirmButtonText: '恢复',
      cancelButtonText: '删除本地内容',
    })
      .then(() => {
        article.value = JSON.parse(localStorage.getItem('write-content')!)
      })
      .catch(() => {
        localStorage.removeItem('write-content')
      })
  }
})
</script>

<style scoped lang="scss">
.img-box {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  border: 1px solid $ac-primary-color;
  border-radius: 5px;

  img {
    width: 100%;
  }

  .el-button {
    position: absolute;
    top: 0;
    right: 0;
    display: none;
  }

  &:hover {
    .el-button {
      display: block;
    }
  }
}
</style>
