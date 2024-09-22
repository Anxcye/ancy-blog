package com.anxcye.service.impl;

import com.anxcye.utils.AdminUtil;
import com.anxcye.utils.SecurityUtil;
import org.springframework.stereotype.Service;

import java.util.Arrays;
import java.util.List;

@Service("ps")
public class PermissionServiceImpl {

    public boolean hasRole(String role) {
        if (AdminUtil.isSuperAdmin(SecurityUtil.getUserId())) {
            return true;
        }

        List<String> roles = SecurityUtil.getLoginUser().getRoles();
        return roles.contains(role);
    }

    public boolean hasPermission(String permission) {
        if (AdminUtil.isSuperAdmin(SecurityUtil.getUserId())) {
            return true;
        }

        List<String> permissions = SecurityUtil.getLoginUser().getPermissions();
        return permissions.contains(permission);
    }

    public boolean hasAnyRole(String... roles) {
        if (AdminUtil.isSuperAdmin(SecurityUtil.getUserId())) {
            return true;
        }

        return Arrays.stream(roles)
                .anyMatch(
                        role -> SecurityUtil.getLoginUser()
                                .getRoles()
                                .contains(role));
    }

    public boolean hasAnyPermission(String... permissions) {
        if (AdminUtil.isSuperAdmin(SecurityUtil.getUserId())) {
            return true;
        }

        return Arrays.stream(permissions)
                .anyMatch(
                        permission -> SecurityUtil.getLoginUser()
                                .getPermissions()
                                .contains(permission));
    }

}
