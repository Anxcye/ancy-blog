package com.anxcye.domain.vo;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class AdminUserVo {
    private String token;
    private List<String> permissions;
    private List<String> role;
    private UserInfoVo userInfoVo;
}
