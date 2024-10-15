package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class ReadPageDto extends PageListDto implements Serializable {
    private String source;
    private String content;
    private String author;
    private Integer addFrom;
}