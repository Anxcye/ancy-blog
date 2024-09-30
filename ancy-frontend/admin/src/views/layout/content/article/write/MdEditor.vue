<template>
  <MdEditor
    :model-value="props.modelValue"
    @update:model-value="handleUpdateModelValue"
    ref="editorRef"
    :toolbars="toolbars"
    :input-box-witdh="inputBoxWitdh"
    @on-upload-img="onUploadImg"
    showCodeRowNumber
    autoDetectCode
  ></MdEditor>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { MdEditor, type ExposeParam, type ToolbarNames } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
import { reqUpload } from '@/api/common'

const props = defineProps(['modelValue'])

const emit = defineEmits(['update:modelValue'])

const handleUpdateModelValue = (v: string) => {
  emit('update:modelValue', v)
}

const initialToolbars: ToolbarNames[] = [
  'underline',
  'sub',
  'sup',
  'image',
  'table',
  'mermaid',
  '=',
  'prettier',
  'pageFullscreen',
  'fullscreen',
  'preview',
  'previewOnly',
  'catalog',
]

const toolbars = ref<ToolbarNames[]>([])
const editorRef = ref<ExposeParam>()

const inputBoxWitdh = ref('50%')

const onUploadImg = async (
  files: File[],
  callback: (urls: string[]) => void,
) => {
  const res = await Promise.all(
    files.map((file) => {
      return reqUpload(file)
    }),
  )

  callback(res.map((item) => item.data))
}

const onResize = () => {
  // mobile
  if (window.innerWidth < 768) {
    toolbars.value = [
      'previewOnly',
      '-',
      ...initialToolbars.filter(
        (item) => !(['preview', 'previewOnly'] as any).includes(item),
      ),
    ]
    inputBoxWitdh.value = '100%'
    editorRef.value?.togglePreview(false)
  } else {
    toolbars.value = initialToolbars
    inputBoxWitdh.value = '50%'
    editorRef.value?.togglePreview(true)
  }
}

onMounted(() => {
  onResize()
  window.addEventListener('resize', onResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
})
</script>

<style scoped lang="scss">
.md-editor-v3 {
  width: 100%;
}
</style>
