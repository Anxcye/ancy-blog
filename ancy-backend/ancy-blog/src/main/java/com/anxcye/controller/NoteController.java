package com.anxcye.controller;

import com.anxcye.domain.dto.PageListDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.NoteService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/note")
public class NoteController {
    @Autowired
    private NoteService noteService;

    @GetMapping("/page")
    public ResponseResult<PageResult> notePage(@ParameterObject PageListDto pageListDto) {
        return ResponseResult.success(noteService.getNotesPage(pageListDto));
    }

}
