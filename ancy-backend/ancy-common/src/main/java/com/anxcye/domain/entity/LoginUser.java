package com.anxcye.domain.entity;

import com.alibaba.fastjson.annotation.JSONField;
import com.anxcye.constants.SystemConstants;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;

import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.Objects;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class LoginUser implements UserDetails {

    private User user;
    private List<String> permissions;
    private List<String> roles;

    public LoginUser(User user) {
        this.user = user;
        this.permissions = List.of();
        this.roles = List.of();
    }

    public LoginUser(User user, List<String> permissions, List<String> roles) {
        this.user = user;
        this.permissions = permissions;
        this.roles = roles;
    }

    


    @JSONField(serialize = false)
    private List<GrantedAuthority> grantedAuthorities;


    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() {
        if (grantedAuthorities == null) {
            grantedAuthorities = new ArrayList<>();
            // 添加角色
            roles.forEach(role -> grantedAuthorities.add(new SimpleGrantedAuthority("ROLE_" + role)));
            // 添加权限
            permissions.forEach(permission -> grantedAuthorities.add(new SimpleGrantedAuthority(permission)));
        }
        return grantedAuthorities;
    }

    @Override
    public String getPassword() {
        return user.getPassword();
    }

    @Override
    public String getUsername() {
        return user.getUserName();
    }

    @Override
    public boolean isAccountNonExpired() {
        return Objects.equals(user.getStatus(), SystemConstants.USER_ENABLE);
    }

    @Override
    public boolean isAccountNonLocked() {
        return Objects.equals(user.getStatus(), SystemConstants.USER_ENABLE);

    }

    @Override
    public boolean isCredentialsNonExpired() {
        return Objects.equals(user.getStatus(), SystemConstants.USER_ENABLE);
    }

    @Override
    public boolean isEnabled() {
        return Objects.equals(user.getStatus(), SystemConstants.USER_ENABLE);
    }
}
