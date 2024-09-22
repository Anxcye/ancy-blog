package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class UserListDto extends PageListDto implements Serializable {
    private String userName;
    private String nickName;
    private String type;
    private String status;
    private String email;
    private String phonenumber;
    private String sex;
}