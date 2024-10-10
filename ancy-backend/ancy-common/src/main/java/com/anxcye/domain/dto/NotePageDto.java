package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class NotePageDto extends PageListDto implements Serializable {
    private String content;
    private String status;
}