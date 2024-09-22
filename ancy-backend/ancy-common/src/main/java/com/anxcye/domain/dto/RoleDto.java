package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;
import java.util.List;

/**
 * 角色信息表
 * @TableName sys_role
 */
@Data
public class RoleDto implements Serializable {
    /**
     * 角色名称
     */
    private String roleName;

    /**
     * 角色权限字符串
     */
    private String roleKey;

    /**
     * 显示顺序
     */
    private Integer roleSort;

    /**
     * menu ids
     */
    private List<Long> menuIds;

    /**
     * 角色状态（0正常 1停用）
     */
    private String status;

    /**
     * 备注
     */
    private String remark;

}