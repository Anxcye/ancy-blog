package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class CommentListDto extends PageListDto implements Serializable {
    private Long articleId;
    
    private String summary;
}