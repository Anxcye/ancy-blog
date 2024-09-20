package com.anxcye.service;

import com.anxcye.domain.entity.Role;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

/**
* @author axy
* @description 针对表【sys_role(角色信息表)】的数据库操作Service
* @createDate 2024-09-20 13:25:43
*/
public interface RoleService extends IService<Role> {

    List<String> getRoleByUserId(Long id);
}
