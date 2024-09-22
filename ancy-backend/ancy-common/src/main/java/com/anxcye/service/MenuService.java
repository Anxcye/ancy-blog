package com.anxcye.service;

import com.anxcye.domain.dto.MenuDto;
import com.anxcye.domain.dto.MenuListDto;
import com.anxcye.domain.entity.Menu;
import com.anxcye.domain.vo.MenuVo;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

/**
* @author axy
* @description 针对表【sys_menu(菜单权限表)】的数据库操作Service
* @createDate 2024-09-20 13:25:34
*/
public interface MenuService extends IService<Menu> {

    boolean hasChild(Long id);

    List<String> getPermissionsByUserId(Long id);

    List<MenuVo> selectMenuTreeByUserId(Long userId);

    List<MenuVo> listMenus(MenuListDto menuListDto);

    boolean addMenu(MenuDto menuDto);

    boolean updateMenu(Long id, MenuDto menuDto);

    boolean deleteMenu(Long id);
}
