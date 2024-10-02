<template>
  <div class="app-container">
    <el-form ref="queryForm" :model="queryParams" size="small" style="margin-bottom: 10px">
      <el-row :gutter="10">
        <el-col :span="10">
          <el-input
            v-model="queryParams.name"
            placeholder="角色名称"
            clearable
            @keyup.enter="getRolePage()"
          />
        </el-col>

        <el-col :span="10">
          <el-select v-model="queryParams.status" placeholder="角色状态" clearable>
            <el-option key="0" label="正常" value="0" />
            <el-option key="1" label="停用" value="1" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" :icon="Search" @click="getRolePage()">搜索</el-button>
        </el-col>
      </el-row>
    </el-form>

    <el-row :gutter="10">
      <el-col :span="1.5">
        <el-button type="primary" plain :icon="Plus" @click="handleAdd">新增</el-button>
      </el-col>
    </el-row>

    <el-table v-loading="loading" :data="roleList">
      <el-table-column label="ID" prop="id" width="50" />
      <el-table-column label="角色名称" prop="roleName" :show-overflow-tooltip="true" width="150" />
      <el-table-column label="权限字符" prop="roleKey" :show-overflow-tooltip="true" width="150" />
      <el-table-column label="显示顺序" prop="roleSort" width="100" />
      <el-table-column label="启用" align="center" width="100">
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
      <el-table-column label="创建时间" align="center" prop="createTime" width="180">
        <template v-slot="scope">
          <span>{{ scope.row.createTime }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" fixed="right" width="170">
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
      @current-change="getRolePage"
      @size-change="getRolePage()"
    />

    <!-- 添加或修改角色配置对话框 -->
    <el-dialog :title="title" v-model="open" class="dialog-form">
      <el-form ref="form" :model="role" :rules="rules" label-width="auto" label-position="top">
        <el-form-item label="角色名称" prop="roleName">
          <el-input v-model="role.roleName" placeholder="角色名称" />
        </el-form-item>

        <el-form-item prop="roleKey">
          <template v-slot:label>
            权限字符
            <span>
              <el-tooltip
                content="控制器中定义的权限字符，如：@PreAuthorize(`@ss.hasRole('admin')`) 中的 admin"
                placement="top"
              >
                <el-icon :size="16">
                  <QuestionFilled />
                </el-icon>
              </el-tooltip>
            </span>
          </template>
          <el-input v-model="role.roleKey" placeholder="权限字符" />
        </el-form-item>

        <el-form-item label="角色顺序" prop="roleSort">
          <el-input-number v-model="role.roleSort" controls-position="right" :min="0" />
        </el-form-item>

        <el-form-item label="状态">
          <el-radio-group v-model="role.status">
            <el-radio key="0" label="0">正常</el-radio>
            <el-radio key="1" label="1">停用</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="菜单权限">
          <el-tree
            ref="menuRef"
            :data="menuOptions"
            show-checkbox
            node-key="id"
            empty-text="加载中，请稍候"
            :props="defaultProps"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="role.remark" type="textarea" placeholder="请输入内容" />
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
import { reqRoleAdd, reqRoleDelete, reqRolePage, reqRoleUpdate } from '@/api/system/role'
import type { RoleAddParams, RoleListData, RolePageParams } from '@/api/system/role/type'
import { ref, onMounted } from 'vue'
import { Plus, Edit, Delete, Search, QuestionFilled } from '@element-plus/icons-vue'
import type { FormRules } from 'element-plus'
import type { MenuListData } from '@/api/system/menu/type'
import { reqMenuTree, reqMenuListByRoleId } from '@/api/system/menu'
const statusLoading = ref<boolean>(false)

const queryParams = ref<RolePageParams>({
  pageNum: 1,
  pageSize: 10,
})
const roleList = ref<RoleListData[]>([])
const total = ref<number>(0)
const loading = ref<boolean>(false)
const open = ref<boolean>(false)
const title = ref<string>('')
const role = ref<RoleAddParams>({})
const rules = ref<FormRules>({})
const menuOptions = ref<MenuListData[]>([])
const defaultProps = ref<object>({
  children: 'children',
  label: 'menuName',
})
const menuRef = ref()

const getRolePage = async (page: number = 1) => {
  loading.value = true
  queryParams.value.pageNum = page
  const res = await reqRolePage(queryParams.value)
  roleList.value = res.data.rows
  total.value = res.data.total
  loading.value = false
}

const getMenuList = async () => {
  const res = await reqMenuTree()
  menuOptions.value = res.data
}

const handleAdd = () => {
  role.value = {}
  menuRef.value?.setCheckedKeys([], true)
  open.value = true
  title.value = '新增'
}
const handleStatusChange = async (row: RoleListData) => {
  statusLoading.value = true
  await reqRoleUpdate(row.id, { status: row.status })
  statusLoading.value = false
}

const handleUpdate = async (row: RoleListData) => {
  role.value = row
  open.value = true
  title.value = '修改' + row.roleName

  const roleMenus = await reqMenuListByRoleId(row.id)
  roleMenus.data.forEach((item) => {
    menuRef.value?.setChecked(item, true, false)
  })
}

const handleDelete = (row: RoleListData) => {
  ElMessageBox.confirm('删除?', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async () => {
    await reqRoleDelete(row.id)
    await getRolePage(queryParams.value.pageNum)
  })
}

const setRoleMenuIds = () => {
  const checkedKeys = menuRef.value?.getCheckedKeys()
  const halfCheckedKeys = menuRef.value?.getHalfCheckedKeys()
  role.value.menuIds = [...checkedKeys, ...halfCheckedKeys]
}

const handleSubmit = async () => {
  if (role.value.id) {
    setRoleMenuIds()
    await reqRoleUpdate(role.value.id, role.value)
  } else {
    setRoleMenuIds()
    await reqRoleAdd(role.value)
    ElMessage.success('新增成功')
  }
  open.value = false
  getRolePage()
}

const handleCancel = () => {
  role.value = {}
  menuRef.value?.setCheckedKeys([], true)
  open.value = false
}

onMounted(() => {
  getRolePage()
  getMenuList()
})
</script>

<style scoped lang="scss"></style>
