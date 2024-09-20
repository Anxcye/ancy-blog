package com.anxcye.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.RoleMenu;
import com.anxcye.service.RoleMenuService;
import com.anxcye.mapper.RoleMenuMapper;
import org.springframework.stereotype.Service;

/**
* @author axy
* @description 针对表【sys_role_menu(角色和菜单关联表)】的数据库操作Service实现
* @createDate 2024-09-20 14:19:33
*/
@Service
public class RoleMenuServiceImpl extends ServiceImpl<RoleMenuMapper, RoleMenu>
    implements RoleMenuService{

}




