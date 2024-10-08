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

        try {
            List<String> roles = SecurityUtil.getLoginUser().getRoles();
            return roles.contains(role);
        } catch (Exception e) {
            return false;
        }
    }

    public boolean hasPermission(String permission) {
        if (AdminUtil.isSuperAdmin(SecurityUtil.getUserId())) {
            return true;
        }
        try {
            List<String> permissions = SecurityUtil.getLoginUser().getPermissions();
            return permissions.contains(permission);
        } catch (Exception e) {
            return false;
        }
    }

    public boolean hasAnyRole(String... roles) {
        if (AdminUtil.isSuperAdmin(SecurityUtil.getUserId())) {
            return true;
        }
        try {

            return Arrays.stream(roles)
                    .anyMatch(
                            role -> SecurityUtil.getLoginUser()
                                    .getRoles()
                                    .contains(role));
        } catch (Exception e) {
            return false;
        }
    }

    public boolean hasAnyPermission(String... permissions) {
        if (AdminUtil.isSuperAdmin(SecurityUtil.getUserId())) {
            return true;
        }
        try {
            return Arrays.stream(permissions)
                    .anyMatch(
                            permission -> SecurityUtil.getLoginUser()
                                    .getPermissions()
                                    .contains(permission));
        } catch (Exception e) {
            return false;
        }
    }

}
