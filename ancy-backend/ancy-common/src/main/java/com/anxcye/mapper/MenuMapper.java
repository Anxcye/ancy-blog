package com.anxcye.mapper;

import com.anxcye.domain.entity.Menu;
import com.baomidou.mybatisplus.core.mapper.BaseMapper;

import java.util.List;

/**
* @author axy
* @description 针对表【sys_menu(菜单权限表)】的数据库操作Mapper
* @createDate 2024-09-20 13:25:34
* @Entity com.anxcye.domain.entity.Menu
*/
public interface MenuMapper extends BaseMapper<Menu> {

    List<Menu> getPermissionsByUserId(Long userId);
}




