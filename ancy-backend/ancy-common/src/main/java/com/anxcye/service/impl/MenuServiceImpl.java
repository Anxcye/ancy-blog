package com.anxcye.service.impl;

import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.entity.Menu;
import com.anxcye.domain.vo.MenuVo;
import com.anxcye.mapper.MenuMapper;
import com.anxcye.service.MenuService;
import com.anxcye.utils.AdminUtil;
import com.anxcye.utils.BeanCopyUtils;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

import java.util.List;

/**
 * @author axy
 * @description 针对表【sys_menu(菜单权限表)】的数据库操作Service实现
 * @createDate 2024-09-20 13:25:34
 */
@Service
public class MenuServiceImpl extends ServiceImpl<MenuMapper, Menu> implements MenuService {

    @Override
    public List<String> getPermissionsByUserId(Long userId) {
        List<Menu> perms = List.of();
        if (AdminUtil.isSuperAdmin(userId)) {
            LambdaQueryWrapper<Menu> menuLambdaQueryWrapper = new LambdaQueryWrapper<>();
            menuLambdaQueryWrapper.in(Menu::getMenuType, SystemConstants.MENU_TABLE_MENU,
                    SystemConstants.MENU_TABLE_BUTTON);
            menuLambdaQueryWrapper.eq(Menu::getStatus, SystemConstants.STATUS_NORMAL);
            menuLambdaQueryWrapper.select(Menu::getPerms);

            perms = list(menuLambdaQueryWrapper);

        } else {
            perms = getBaseMapper().getPermissionsByUserId(userId);
        }
        return perms.stream().map(Menu::getPerms).toList();

    }

    private List<MenuVo> buildMenuTree(List<MenuVo> allMenus, Long parentId) {
        return allMenus.stream()
                .filter(menuVo -> menuVo.getParentId().equals(parentId))
                .peek(menuVo -> menuVo.setChildren(buildMenuTree(allMenus, menuVo.getId())))
                .toList();
    }

    @Override
    public List<MenuVo> selectMenuTreeByUserId(Long userId) {
        List<Menu> menus;
        if (AdminUtil.isSuperAdmin(userId)) {
            LambdaQueryWrapper<Menu> menuLambdaQueryWrapper = new LambdaQueryWrapper<>();
            menuLambdaQueryWrapper.in(Menu::getMenuType, SystemConstants.MENU_TABLE_MENU,
                    SystemConstants.MENU_TABLE_CATALOG);
            menuLambdaQueryWrapper.eq(Menu::getStatus, SystemConstants.STATUS_NORMAL);
            menuLambdaQueryWrapper.orderByAsc(Menu::getParentId, Menu::getOrderNum);
            menus = list(menuLambdaQueryWrapper);
        } else {
            menus = getBaseMapper().selectRoutersByUserId(userId);
        }

        List<MenuVo> menuVos = BeanCopyUtils.copyList(menus, MenuVo.class);
        List<MenuVo> menuVos1 = buildMenuTree(menuVos, SystemConstants.ROOT_MENU_PARENT_ID);
        return menuVos1;
    }
}
