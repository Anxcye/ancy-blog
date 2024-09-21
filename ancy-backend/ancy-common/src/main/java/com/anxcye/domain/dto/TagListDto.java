package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class TagListDto implements Serializable {
    private Integer pageNum;
    private Integer pageSize;
    private String name;
    private String remark;



}