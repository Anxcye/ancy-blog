package com.anxcye.service;

import com.anxcye.domain.entity.Menu;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

/**
* @author axy
* @description 针对表【sys_menu(菜单权限表)】的数据库操作Service
* @createDate 2024-09-20 13:25:34
*/
public interface MenuService extends IService<Menu> {

    List<String> getPermissionsByUserId(Long id);
}
