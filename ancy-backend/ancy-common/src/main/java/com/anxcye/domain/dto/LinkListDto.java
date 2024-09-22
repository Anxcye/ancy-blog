package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class LinkListDto extends PageListDto implements Serializable {
    private String title;
    private String status;
}