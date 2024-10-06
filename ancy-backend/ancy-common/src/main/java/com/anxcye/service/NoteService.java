package com.anxcye.service;

import com.anxcye.domain.dto.PageListDto;
import com.anxcye.domain.entity.Note;
import com.anxcye.domain.result.PageResult;
import com.baomidou.mybatisplus.extension.service.IService;

/**
* @author axy
* @description 针对表【ancy_note(note表)】的数据库操作Service
* @createDate 2024-10-06 19:44:06
*/
public interface NoteService extends IService<Note> {

    PageResult getNotesPage(PageListDto pageListDto);
}
