package com.anxcye.service.impl;

import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.PageListDto;
import com.anxcye.domain.entity.Note;
import com.anxcye.domain.result.PageResult;
import com.anxcye.mapper.NoteMapper;
import com.anxcye.service.NoteService;
import com.anxcye.utils.SecurityUtil;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

/**
 * @author axy
 * @description 针对表【ancy_note(note表)】的数据库操作Service实现
 * @createDate 2024-10-06 19:44:06
 */
@Service
public class NoteServiceImpl extends ServiceImpl<NoteMapper, Note>
        implements NoteService {

    private LambdaQueryWrapper<Note> getNoteWrapper() {
        LambdaQueryWrapper<Note> wrapper = new LambdaQueryWrapper<>();
        if (!SecurityUtil.isAdmin()) {
            wrapper.eq(Note::getStatus, SystemConstants.STATUS_NORMAL);
        }
        return wrapper;
    }

    @Override
    public PageResult getNotesPage(PageListDto pageListDto) {
        LambdaQueryWrapper<Note> wrapper = getNoteWrapper();
        wrapper.orderByDesc(Note::getIsTop)
                .orderByAsc(Note::getOrderNum)
                .orderByDesc(Note::getCreateTime);
        Page<Note> page = new Page<>(pageListDto.getPageNum(), pageListDto.getPageSize());
        page(page, wrapper);
        for (Note note : page.getRecords()) {
            note.setViewCount(note.getViewCount() + 1);
            updateById(note);
        }
        return new PageResult(page.getTotal(), page.getRecords());
    }
}
