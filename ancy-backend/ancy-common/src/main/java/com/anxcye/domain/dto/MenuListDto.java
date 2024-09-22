package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class MenuListDto implements Serializable {
    private String name;
    private String status;
}