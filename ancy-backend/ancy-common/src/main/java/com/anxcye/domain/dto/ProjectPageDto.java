package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class ProjectPageDto extends PageListDto implements Serializable {
    private String title;
    private String summary;
    private String status;
    private String type;
}