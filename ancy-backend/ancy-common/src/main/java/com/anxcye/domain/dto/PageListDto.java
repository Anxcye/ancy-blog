package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class PageListDto implements Serializable {
    private Integer pageNum;
    private Integer pageSize;
}