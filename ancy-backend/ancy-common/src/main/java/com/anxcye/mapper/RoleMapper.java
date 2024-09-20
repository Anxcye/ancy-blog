package com.anxcye.mapper;

import com.anxcye.domain.entity.Role;
import com.baomidou.mybatisplus.core.mapper.BaseMapper;

import java.util.List;

/**
* @author axy
* @description 针对表【sys_role(角色信息表)】的数据库操作Mapper
* @createDate 2024-09-20 13:25:43
* @Entity com.anxcye.domain.entity.Role
*/
public interface RoleMapper extends BaseMapper<Role> {

    List<Role> getRolesByUserId(Long userId);
}




