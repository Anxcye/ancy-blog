package com.anxcye.controller;

import com.anxcye.domain.dto.NoteDto;
import com.anxcye.domain.dto.NotePageDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.NoteVo;
import com.anxcye.service.NoteService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/note")
public class NoteController {

    @Autowired
    private NoteService noteService;

    @PreAuthorize("@ps.hasPermission('content:note:list')")
    @GetMapping
    public ResponseResult<PageResult> notePage(@ParameterObject NotePageDto notePageDto) {
        return ResponseResult.success(noteService.getNotePage(notePageDto, true));
    }

    @PreAuthorize("@ps.hasPermission('content:note:query')")
    @GetMapping("/{id}")
    public ResponseResult<NoteVo> noteGetById(@PathVariable Long id) {
        return ResponseResult.success(noteService.getNoteById(id));
    }

    @PreAuthorize("@ps.hasPermission('content:note:add')")
    @PostMapping
    public ResponseResult<Long> noteAdd(@RequestBody NoteDto noteDto) {
        return ResponseResult.success(noteService.addNote(noteDto));
    }

    @PreAuthorize("@ps.hasPermission('content:note:edit')")
    @PutMapping("/{id}")
    public ResponseResult<Boolean> noteUpdate(@PathVariable Long id, @RequestBody NoteDto noteDto) {
        return ResponseResult.success(noteService.updateNote(id, noteDto));
    }

    @PreAuthorize("@ps.hasPermission('content:note:remove')")
    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> noteDelete(@PathVariable Long id) {
        return ResponseResult.success(noteService.deleteNote(id));
    }
}
