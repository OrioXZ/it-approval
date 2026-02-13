import { TestBed } from '@angular/core/testing';

import { It03 } from './it03';

describe('It03', () => {
  let service: It03;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(It03);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
