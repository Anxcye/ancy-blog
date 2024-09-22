package com.anxcye.service.impl;

import com.anxcye.domain.dto.RoleDto;
import com.anxcye.domain.dto.RoleListDto;
import com.anxcye.domain.entity.Role;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.RoleVo;
import com.anxcye.mapper.RoleMapper;
import com.anxcye.service.RoleMenuService;
import com.anxcye.service.RoleService;
import com.anxcye.utils.BeanCopyUtils;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

import java.util.List;

/**
 * @author axy
 * @description 针对表【sys_role(角色信息表)】的数据库操作Service实现
 * @createDate 2024-09-20 13:25:43
 */
@Service
public class RoleServiceImpl extends ServiceImpl<RoleMapper, Role>
        implements RoleService {

    @Autowired
    private RoleMenuService roleMenuService;

    @Override
    public List<String> getRoleByUserId(Long userId) {
        return getBaseMapper().getRolesByUserId(userId).stream()
                .map(Role::getRoleKey).toList();

    }

    @Override
    public PageResult getPage(RoleListDto roleListDto) {
        LambdaQueryWrapper<Role> wrapper = new LambdaQueryWrapper<>();

        wrapper.like(StringUtils.hasText(roleListDto.getName()), Role::getRoleName, roleListDto.getName());
        wrapper.eq(StringUtils.hasText(roleListDto.getStatus()), Role::getStatus, roleListDto.getStatus());
        wrapper.orderByAsc(Role::getRoleSort);

        Page<Role> page = new Page<>(roleListDto.getPageNum(), roleListDto.getPageSize());
        page(page, wrapper);

        List<RoleVo> roleVos = BeanCopyUtils.copyList(page.getRecords(), RoleVo.class);
        return new PageResult(page.getTotal(), roleVos);

    }

    @Override
    @Transactional
    public boolean addRole(RoleDto roleDto) {
        Role role = BeanCopyUtils.copyBean(roleDto, Role.class);
        save(role);
        roleMenuService.saveByRoleId(role.getId(), roleDto.getMenuIds());
        return true;
    }

    @Override
    public boolean updateRole(Long id, RoleDto roleDto) {
        Role role = BeanCopyUtils.copyBean(roleDto, Role.class);
        role.setId(id);
        updateById(role);
        roleMenuService.updateByRoleId(id, roleDto.getMenuIds());
        return true;
    }

    @Override
    public boolean deleteRole(Long id) {
        removeById(id);
        return true;
    }

    @Override
    public RoleVo getRoleById(Long id) {
        Role role = getById(id);
        RoleVo roleVo = BeanCopyUtils.copyBean(role, RoleVo.class);
        roleVo.setMenuIds(roleMenuService.getMenuIdsByRoleId(id));
        return roleVo;
    }
}




