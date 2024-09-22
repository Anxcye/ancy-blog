package com.anxcye.service;

import com.anxcye.domain.dto.RoleDto;
import com.anxcye.domain.dto.RoleListDto;
import com.anxcye.domain.entity.Role;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.RoleVo;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

/**
* @author axy
* @description 针对表【sys_role(角色信息表)】的数据库操作Service
* @createDate 2024-09-20 13:25:43
*/
public interface RoleService extends IService<Role> {

    List<String> getRoleByUserId(Long id);

    PageResult getPage(RoleListDto roleListDto);

    boolean addRole(RoleDto roleDto);

    boolean updateRole(Long id, RoleDto roleDto);

    boolean deleteRole(Long id);

    RoleVo getRoleById(Long id);
}
