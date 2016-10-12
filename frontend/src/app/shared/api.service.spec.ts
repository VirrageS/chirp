import { inject, async, TestBed } from '@angular/core/testing';
import { BaseRequestOptions, Http, Response, ResponseOptions } from '@angular/http';
import { MockBackend } from '@angular/http/testing';

import { ApiService } from './api.service';


describe('ApiSerivce', () => {
  let apiSerivce: ApiService;
  let mockbackend: MockBackend;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        BaseRequestOptions,
        MockBackend,
        {
          provide: Http,
          useFactory: (backend, options) => new Http(backend, options),
          deps: [MockBackend, BaseRequestOptions]
        },
        ApiService
      ]
    })
  })

  beforeEach(inject([ApiService, MockBackend], (service, mock) => {
    apiSerivce = service;
    mockbackend = mock;
  }))

  it('should make get request', () => {
     let connection;
     let response = [
       {title: "Title", description: "", url: "github.com/someth/ing", points: 20},
       {title: "Title", description: "", url: "github.com/someth/ing", points: 20},
       {title: "Title", description: "", url: "github.com/someth/ing", points: 20}
     ];

     mockbackend.connections.subscribe(connection => {
       connection.mockRespond(new Response(
         new ResponseOptions({body: JSON.stringify(response), status: 200})
       ))
     })

     apiSerivce.get('/github')
       .subscribe(notes => {
         expect(notes).toEqual(response);
       })
  })
})
