package com.anxcye.service.impl;

import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.MenuDto;
import com.anxcye.domain.dto.MenuListDto;
import com.anxcye.domain.entity.Menu;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.vo.MenuVo;
import com.anxcye.exception.SystemException;
import com.anxcye.mapper.MenuMapper;
import com.anxcye.service.MenuService;
import com.anxcye.utils.AdminUtil;
import com.anxcye.utils.BeanCopyUtils;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.util.List;
import java.util.Objects;


/**
 * @author axy
 * @description 针对表【sys_menu(菜单权限表)】的数据库操作Service实现
 * @createDate 2024-09-20 13:25:34
 */
@Service
public class MenuServiceImpl extends ServiceImpl<MenuMapper, Menu> implements MenuService {

    private List<MenuVo> buildMenuTree(List<MenuVo> allMenus, Long parentId) {
        return allMenus.stream()
                .filter(menuVo -> menuVo.getParentId().equals(parentId))
                .peek(menuVo -> menuVo.setChildren(buildMenuTree(allMenus, menuVo.getId())))
                .toList();
    }

    @Override
    public boolean hasChild(Long id) {
        LambdaQueryWrapper<Menu> wrapper = new LambdaQueryWrapper<>();
        wrapper.eq(Menu::getParentId, id);
        return count(wrapper) > 0;
    }

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
        return buildMenuTree(menuVos, SystemConstants.ROOT_MENU_PARENT_ID);
    }

    @Override
    public List<MenuVo> listMenus(MenuListDto menuListDto) {
        LambdaQueryWrapper<Menu> wrapper = new LambdaQueryWrapper<>();

        wrapper.like(StringUtils.hasText(menuListDto.getName()), Menu::getMenuName, menuListDto.getName());
        wrapper.eq(StringUtils.hasText(menuListDto.getStatus()), Menu::getStatus, menuListDto.getStatus());

        wrapper.orderByAsc(Menu::getParentId, Menu::getOrderNum);
        List<Menu> menus = list(wrapper);
        return BeanCopyUtils.copyList(menus, MenuVo.class);
    }

    @Override
    public boolean addMenu(MenuDto menuDto) {
        Menu menu = BeanCopyUtils.copyBean(menuDto, Menu.class);
        save(menu);
        return true;
    }

    @Override
    public boolean updateMenu(Long id, MenuDto menuDto) {
        if (Objects.equals(menuDto.getParentId(), id)) {
            throw new SystemException(AppHttpCodeEnum.SELF_PARENT_ERROR);
        }
        Menu menu = BeanCopyUtils.copyBean(menuDto, Menu.class);
        menu.setId(id);
        updateById(menu);
        return true;
    }

    @Override
    public boolean deleteMenu(Long id) {
        if (hasChild(id)) {
            throw new SystemException(AppHttpCodeEnum.HAS_CHILD_DELETE_FAILED);
        }
        removeById(id);
        return true;
    }

}
