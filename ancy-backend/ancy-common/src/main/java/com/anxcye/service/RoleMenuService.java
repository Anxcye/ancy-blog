package com.anxcye.service;

import com.anxcye.domain.entity.RoleMenu;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

/**
* @author axy
* @description 针对表【sys_role_menu(角色和菜单关联表)】的数据库操作Service
* @createDate 2024-09-20 14:19:33
*/
public interface RoleMenuService extends IService<RoleMenu> {

    boolean saveByRoleId(Long RoleId, List<Long> menuIds);

    boolean deleteByRoleId(Long roleId);

    boolean updateByRoleId(Long roleId, List<Long> menuIds);

    List<Long> getMenuIdsByRoleId(Long roleId);

    boolean deleteByMenuId(Long menuId);
}
