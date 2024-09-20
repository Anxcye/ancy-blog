package com.anxcye.service.impl;

import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.entity.Menu;
import com.anxcye.mapper.MenuMapper;
import com.anxcye.service.MenuService;
import com.anxcye.utils.AdminUtil;
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
            menuLambdaQueryWrapper.in(Menu::getMenuType, SystemConstants.MENU_TABLE_MENU, SystemConstants.MENU_TABLE_BUTTON);
            menuLambdaQueryWrapper.eq(Menu::getStatus, SystemConstants.STATUS_NORMAL);
            menuLambdaQueryWrapper.select(Menu::getPerms);

            perms = list(menuLambdaQueryWrapper);

        } else {
            perms = getBaseMapper().getPermissionsByUserId(userId);
        }
        return perms.stream()
                .map(Menu::getPerms)
                .toList();

    }
}




