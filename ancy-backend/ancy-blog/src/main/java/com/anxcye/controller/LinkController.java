package com.anxcye.controller;

import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.LinkService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/link")
public class LinkController {

    @Autowired
    private LinkService linkService;

    @GetMapping("/list")
    public ResponseResult list() {
        return ResponseResult.okResult(linkService.getApprovedLinks());
    }

}
