package com.anxcye.service.impl;

import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.NoteDto;
import com.anxcye.domain.dto.NotePageDto;
import com.anxcye.domain.entity.Note;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.NoteVo;
import com.anxcye.mapper.NoteMapper;
import com.anxcye.service.NoteService;
import com.anxcye.utils.BeanCopyUtils;
import com.anxcye.utils.SecurityUtil;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

import java.util.List;

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

//    @Override
//    public PageResult getNotesPage(PageListDto pageListDto) {
//        LambdaQueryWrapper<Note> wrapper = getNoteWrapper();
//        wrapper.orderByDesc(Note::getIsTop)
//                .orderByAsc(Note::getOrderNum)
//                .orderByDesc(Note::getCreateTime);
//        Page<Note> page = new Page<>(pageListDto.getPageNum(), pageListDto.getPageSize());
//        page(page, wrapper);
//        for (Note note : page.getRecords()) {
//            note.setViewCount(note.getViewCount() + 1);
//            updateById(note);
//        }
//        List<NoteVo> noteVos = BeanCopyUtils.copyList(page.getRecords(), NoteVo.class);
//        return new PageResult(page.getTotal(), noteVos);
//    }

    @Override
    public PageResult getNotePage(NotePageDto notePageDto, Boolean isAdmin) {
        LambdaQueryWrapper<Note> wrapper = getNoteWrapper();
        wrapper.like(notePageDto.getContent()!= null, Note::getContent, notePageDto.getContent())
             .eq(notePageDto.getStatus()!= null, Note::getStatus, notePageDto.getStatus())
             .orderByDesc(Note::getIsTop)
             .orderByAsc(Note::getOrderNum)
             .orderByDesc(Note::getCreateTime);
        Page<Note> page = new Page<>(notePageDto.getPageNum(), notePageDto.getPageSize());
        page(page, wrapper);
        if (!isAdmin){
            for (Note note : page.getRecords()) {
                note.setViewCount(note.getViewCount() + 1);
                updateById(note);
            }
        }
        List<NoteVo> noteVos = BeanCopyUtils.copyList(page.getRecords(), NoteVo.class);
        return new PageResult(page.getTotal(), noteVos);
    }

    @Override
    public NoteVo getNoteById(Long id) {
        LambdaQueryWrapper<Note> wrapper = getNoteWrapper();
        wrapper.eq(Note::getId, id);
        Note note = getOne(wrapper);
        return BeanCopyUtils.copyBean(note, NoteVo.class);
    }

    @Override
    public Long addNote(NoteDto noteDto) {
        Note note = BeanCopyUtils.copyBean(noteDto, Note.class);
        save(note);
        return note.getId();
    }

    @Override
    public Boolean updateNote(Long id, NoteDto noteDto) {
        Note note = BeanCopyUtils.copyBean(noteDto, Note.class);
        note.setId(id);
        return updateById(note);
    }

    @Override
    public Boolean deleteNote(Long id) {
        return removeById(id);
    }
}
